package lab

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReviewCanonicalizesVerdictForScoring(t *testing.T) {
	service := NewService(NewMemoryStore(), DefaultCatalog())
	SeedDemo(service)
	server := NewServer(service)

	reviewRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/observations/demo-001/review",
		strings.NewReader(`{"verdict":" Confirmed ","confidence":0.9,"notes":"Independent audio confirms the field note."}`),
	)
	reviewRecorder := httptest.NewRecorder()
	server.ServeHTTP(reviewRecorder, reviewRequest)
	if reviewRecorder.Code != http.StatusOK {
		t.Fatalf("review status=%d", reviewRecorder.Code)
	}

	var item Observation
	if err := json.NewDecoder(reviewRecorder.Body).Decode(&item); err != nil {
		t.Fatal(err)
	}
	if item.Review == nil || item.Review.Verdict != "confirmed" {
		t.Fatalf("review=%+v", item.Review)
	}

	reportRequest := httptest.NewRequest(http.MethodGet, "/api/v1/observations/demo-001/report", nil)
	reportRecorder := httptest.NewRecorder()
	server.ServeHTTP(reportRecorder, reportRequest)
	if reportRecorder.Code != http.StatusOK {
		t.Fatalf("report status=%d", reportRecorder.Code)
	}
	var report SignalReport
	if err := json.NewDecoder(reportRecorder.Body).Decode(&report); err != nil {
		t.Fatal(err)
	}
	if report.Score != 82 {
		t.Fatalf("score=%d, want 82", report.Score)
	}
}
