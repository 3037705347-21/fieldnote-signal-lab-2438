package lab

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddingLabelsPreservesReviewedState(t *testing.T) {
	service := NewService(NewMemoryStore(), DefaultCatalog())
	SeedDemo(service)
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/observations/demo-001/labels",
		strings.NewReader(`{"labels":["weather"]}`),
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
	if item.State != StateReviewed {
		t.Fatalf("state=%s", item.State)
	}
}
