package lab

import (
	"sync"
	"testing"
	"time"
)

type labelBarrierStore struct {
	inner   Store
	target  string
	mu      sync.Mutex
	gets    int
	release chan struct{}
}

func (s *labelBarrierStore) Create(item Observation) error { return s.inner.Create(item) }
func (s *labelBarrierStore) Update(item Observation) error  { return s.inner.Update(item) }
func (s *labelBarrierStore) AddLabels(id string, labels []string) error {
	return s.inner.AddLabels(id, labels)
}
func (s *labelBarrierStore) List() []Observation            { return s.inner.List() }
func (s *labelBarrierStore) Get(id string) (Observation, error) {
	item, err := s.inner.Get(id)
	if err != nil || id != s.target {
		return item, err
	}
	s.mu.Lock()
	s.gets++
	block := s.gets <= 2
	if s.gets == 2 {
		close(s.release)
	}
	s.mu.Unlock()
	if block {
		<-s.release
	}
	return item, nil
}

func TestConcurrentLabelsAreMerged(t *testing.T) {
	base := NewMemoryStore()
	creator := NewService(base, DefaultCatalog())
	item, err := creator.Create(CreateObservationInput{
		Site:        "wetland",
		ObservedAt:  time.Now().UTC().Add(-time.Hour).Format(time.RFC3339),
		Description: "Research team recorded repeated calls beside the marsh reeds.",
	})
	if err != nil {
		t.Fatal(err)
	}
	store := &labelBarrierStore{inner: base, target: item.ID, release: make(chan struct{})}
	service := NewService(store, DefaultCatalog())

	var wg sync.WaitGroup
	errs := make(chan error, 2)
	start := make(chan struct{})
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		_, callErr := service.AddLabels(item.ID, AddLabelsInput{Labels: []string{"acoustic"}})
		errs <- callErr
	}()
	go func() {
		defer wg.Done()
		<-start
		_, callErr := service.AddLabels(item.ID, AddLabelsInput{Labels: []string{"nesting"}})
		errs <- callErr
	}()
	close(start)
	wg.Wait()
	close(errs)
	for callErr := range errs {
		if callErr != nil {
			t.Fatal(callErr)
		}
	}
	stored, err := service.Get(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !hasLabel(stored, "acoustic") || !hasLabel(stored, "nesting") {
		t.Fatalf("labels=%v", stored.Labels)
	}
}
