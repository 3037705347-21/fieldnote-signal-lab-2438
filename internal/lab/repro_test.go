package lab

import (
	"testing"
	"time"
)

func TestObservationReadsAreIsolated(t *testing.T) {
	service := NewService(NewMemoryStore(), DefaultCatalog())
	item, err := service.Create(CreateObservationInput{
		Site:        "wetland",
		ObservedAt:  time.Now().UTC().Add(-time.Hour).Format(time.RFC3339),
		Description: "Research team recorded repeated calls beside the marsh reeds.",
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.AddLabels(item.ID, AddLabelsInput{Labels: []string{"acoustic", "nesting"}}); err != nil {
		t.Fatal(err)
	}

	read, err := service.Get(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	read.Labels[0] = "pollinator"

	stored, err := service.Get(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !hasLabel(stored, "acoustic") || hasLabel(stored, "pollinator") {
		t.Fatalf("stored labels were mutated through a read: %v", stored.Labels)
	}
}
