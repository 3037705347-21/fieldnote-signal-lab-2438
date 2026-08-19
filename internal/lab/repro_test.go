package lab

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReviewNormalizesVerdict(t *testing.T) {
	service := NewService(NewMemoryStore(), DefaultCatalog())
	SeedDemo(service)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/observations/demo-001/review",
		strings.NewReader(`{"verdict":"  CONFIRMED  ","confidence":0.9,"notes":"Independent recording confirms the observation."}`),
	)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	NewServer(service).ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}

	item, err := service.Get("demo-001")
	if err != nil {
		t.Fatal(err)
	}
	if item.State != StateReviewed || item.Review == nil {
		t.Fatalf("item=%+v", item)
	}
	if item.Review.Verdict != "confirmed" {
		t.Fatalf("verdict=%q", item.Review.Verdict)
	}
}
