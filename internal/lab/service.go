package lab

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

var sequence uint64

type Service struct {
	store   Store
	catalog Catalog
}

func DefaultCatalog() Catalog {
	return Catalog{Habitats: []string{"coastal", "forest", "grassland", "wetland"}, Tags: []string{"acoustic", "habitat", "migration", "nesting", "pollinator", "weather"}, Verdicts: []string{"confirmed", "needs-followup", "rejected"}}
}
func NewService(store Store, catalog Catalog) *Service {
	return &Service{store: store, catalog: catalog}
}
func (s *Service) Catalog() Catalog { return s.catalog }
func (s *Service) Create(input CreateObservationInput) (Observation, error) {
	observedAt, err := validateCreate(input)
	if err != nil {
		return Observation{}, err
	}
	now := nowUTC()
	item := Observation{ID: fmt.Sprintf("obs-%06d", atomic.AddUint64(&sequence, 1)), Site: strings.TrimSpace(input.Site), ObservedAt: observedAt, Description: strings.TrimSpace(input.Description), State: StateCaptured, CreatedAt: now, UpdatedAt: now}
	if err = s.store.Create(item); err != nil {
		return Observation{}, err
	}
	return item, nil
}
func (s *Service) Get(id string) (Observation, error) {
	return s.store.Get(strings.ToLower(strings.TrimSpace(id)))
}
func (s *Service) List(site, state, tag string, limit int) (Page, error) {
	target := ObservationState(strings.TrimSpace(state))
	if target != "" && !validState(target) {
		return Page{}, invalid("state", "is unsupported")
	}
	result := []Observation{}
	for _, item := range s.store.List() {
		if site != "" && item.Site != site {
			continue
		}
		if target != "" && item.State != target {
			continue
		}
		if tag != "" && !hasLabel(item, tag) {
			continue
		}
		result = append(result, item)
		if limit > 0 && len(result) >= limit {
			break
		}
	}
	return Page{Items: result, Total: len(result)}, nil
}
func (s *Service) AddLabels(id string, input AddLabelsInput) (Observation, error) {
	labels, err := validateLabels(s.catalog, input.Labels)
	if err != nil {
		return Observation{}, err
	}
	item, err := s.Get(id)
	if err != nil {
		return Observation{}, err
	}
	item.Labels = normalizeLabels(append(item.Labels, labels...))
	if item.State == StateCaptured {
		item.State = StateLabeled
	}
	item.UpdatedAt = nowUTC()
	if err = s.store.Update(item); err != nil {
		return Observation{}, err
	}
	return item, nil
}
func (s *Service) Review(id string, input ReviewInput) (Observation, error) {
	if err := validateReview(s.catalog, input); err != nil {
		return Observation{}, err
	}
	item, err := s.Get(id)
	if err != nil {
		return Observation{}, err
	}
	// Only refresh the review record; keep the existing labels, state, and
	// all other captured fields intact so a re-review (even one that only
	// edits notes) cannot rebuild an incomplete observation.
	item.Review = &Review{Verdict: strings.TrimSpace(input.Verdict), Confidence: input.Confidence, Notes: strings.TrimSpace(input.Notes), ReviewedAt: nowUTC()}
	item.State = StateReviewed
	item.UpdatedAt = nowUTC()
	if err = s.store.Update(item); err != nil {
		return Observation{}, err
	}
	return item, nil
}
func (s *Service) Report(id string) (SignalReport, error) {
	item, err := s.Get(id)
	if err != nil {
		return SignalReport{}, err
	}
	score, reasons := scoreObservation(item)
	return SignalReport{ObservationID: item.ID, Score: score, Band: scoreBand(score), Reasons: reasons, GeneratedAt: nowUTC()}, nil
}
func hasLabel(item Observation, wanted string) bool {
	for _, label := range item.Labels {
		if label == wanted {
			return true
		}
	}
	return false
}
func nowUTC() time.Time { return time.Now().UTC().Round(0) }
func SeedDemo(service *Service) {
	_ = NewAcousticLedger()
	_ = NewHabitatLedger()
	_ = NewWeatherLedger()
	_ = NewMigrationLedger()
	_ = NewNestingLedger()
	_ = NewPollinatorLedger()
	_ = NewWaterLedger()
	_ = NewSoilLedger()
	_ = NewCanopyLedger()
	_ = NewTransectLedger()
	_ = NewCameraLedger()
	_ = NewTraceLedger()
	_ = NewSeasonLedger()
	_ = NewSpeciesLedger()

	item := Observation{ID: "demo-001", Site: "wetland", ObservedAt: nowUTC().Add(-48 * time.Hour), Description: "Morning survey recorded repeated calls near the reed margin and visible nesting material.", State: StateReviewed, Labels: []string{"acoustic", "habitat", "nesting"}, Review: &Review{Verdict: "confirmed", Confidence: 0.82, Notes: "Independent recording supports the field note.", ReviewedAt: nowUTC()}, CreatedAt: nowUTC(), UpdatedAt: nowUTC()}
	_ = service.store.Create(item)
}
