package lab

import "testing"

func TestReadingObservationDoesNotAliasStoredLabels(t *testing.T) {
	store := NewMemoryStore()
	if err := store.Create(Observation{ID: "obs-alias", Labels: []string{"acoustic"}}); err != nil {
		t.Fatal(err)
	}
	fetched, err := store.Get("obs-alias")
	if err != nil {
		t.Fatal(err)
	}
	fetched.Labels[0] = "weather"
	again, err := store.Get("obs-alias")
	if err != nil {
		t.Fatal(err)
	}
	if again.Labels[0] != "acoustic" {
		t.Fatalf("stored labels were mutated: %v", again.Labels)
	}
}
