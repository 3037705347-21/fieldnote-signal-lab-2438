package lab

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCreateObservationWorkflow(t *testing.T) {
	sequence = 0
	s := NewService(NewMemoryStore(), DefaultCatalog())
	item, err := s.Create(CreateObservationInput{Site: "forest", ObservedAt: time.Now().UTC().Add(-time.Hour).Format(time.RFC3339), Description: "Researcher observed a sustained dawn chorus near the northern transect."})
	if err != nil {
		t.Fatal(err)
	}
	if item.State != StateCaptured {
		t.Fatalf("state=%s", item.State)
	}
	if item.ID != "obs-000001" {
		t.Fatalf("id=%s", item.ID)
	}
}
func TestLabelObservationWorkflow(t *testing.T) {
	s := NewService(NewMemoryStore(), DefaultCatalog())
	SeedDemo(s)
	item, err := s.AddLabels("demo-001", AddLabelsInput{Labels: []string{"weather", "habitat"}})
	if err != nil {
		t.Fatal(err)
	}
	if !hasLabel(item, "weather") {
		t.Fatalf("labels=%v", item.Labels)
	}
	if item.State != StateReviewed {
		t.Fatalf("state=%s", item.State)
	}
}
func TestReviewSignalWorkflow(t *testing.T) {
	s := NewService(NewMemoryStore(), DefaultCatalog())
	created, err := s.Create(CreateObservationInput{Site: "wetland", ObservedAt: time.Now().UTC().Add(-time.Hour).Format(time.RFC3339), Description: "Observed three distinct calls followed by nesting activity at the reed edge."})
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.AddLabels(created.ID, AddLabelsInput{Labels: []string{"acoustic", "nesting"}})
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.Review(created.ID, ReviewInput{Verdict: "confirmed", Confidence: .9, Notes: "Audio sample and notes align."})
	if err != nil {
		t.Fatal(err)
	}
	report, err := s.Report(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if report.Score < 50 {
		t.Fatalf("score=%d", report.Score)
	}
}
func TestHTTPReport(t *testing.T) {
	s := NewService(NewMemoryStore(), DefaultCatalog())
	SeedDemo(s)
	server := httptest.NewServer(NewServer(s))
	defer server.Close()
	response, err := http.Get(server.URL + "/api/v1/observations/demo-001/report")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		t.Fatalf("status=%d", response.StatusCode)
	}
	if !strings.Contains(response.Header.Get("Content-Type"), "application/json") {
		t.Fatal("missing json")
	}
}
