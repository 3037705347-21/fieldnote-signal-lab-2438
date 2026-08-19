package lab

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReviewPreservesExistingLabels(t *testing.T) {
	service := NewService(NewMemoryStore(), DefaultCatalog())
	SeedDemo(service)
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/observations/demo-001/review",
		strings.NewReader(`{"verdict":"confirmed","confidence":0.9,"notes":"Independent recording confirms the observation."}`),
	)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	NewServer(service).ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	var item Observation
	if err := json.NewDecoder(response.Body).Decode(&item); err != nil {
		t.Fatal(err)
	}
	if len(item.Labels) != 3 || !hasLabel(item, "acoustic") || !hasLabel(item, "nesting") {
		t.Fatalf("labels=%v", item.Labels)
	}
}
