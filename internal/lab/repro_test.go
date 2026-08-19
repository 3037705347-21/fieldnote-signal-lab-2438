package lab

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRejectedReviewDoesNotIncreaseReportScore(t *testing.T) {
	sequence = 0
	service := NewService(NewMemoryStore(), DefaultCatalog())
	item, err := service.Create(CreateObservationInput{Site: "wetland", ObservedAt: time.Now().UTC().Add(-time.Hour).Format(time.RFC3339), Description: "Researchers documented steady calls beside a sheltered reed bed before sunrise."})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = service.AddLabels(item.ID, AddLabelsInput{Labels: []string{"acoustic", "nesting"}}); err != nil {
		t.Fatal(err)
	}
	if _, err = service.Review(item.ID, ReviewInput{Verdict: "rejected", Confidence: .9, Notes: "Evidence did not match the field record."}); err != nil {
		t.Fatal(err)
	}
	handler := NewServer(service)
	route := strings.Replace("/api/v1/observations/{id}/report", "{id}", item.ID, 1)
	req := httptest.NewRequest(http.MethodGet, route, nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var report SignalReport
	if err := json.NewDecoder(rr.Body).Decode(&report); err != nil {
		t.Fatal(err)
	}
	if report.Score != 49 {
		t.Fatalf("score=%d reasons=%v", report.Score, report.Reasons)
	}
}
