package lab

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddingLabelKeepsReviewedObservationReviewed(t *testing.T) {
	service := NewService(NewMemoryStore(), DefaultCatalog())
	SeedDemo(service)
	handler := NewServer(service)
	route := strings.Replace("/api/v1/observations/{id}/labels", "{id}", "demo-001", 1)
	req := httptest.NewRequest(http.MethodPost, route, strings.NewReader(`{"labels":["weather"]}`))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var item Observation
	if err := json.NewDecoder(rr.Body).Decode(&item); err != nil {
		t.Fatal(err)
	}
	if item.State != StateReviewed {
		t.Fatalf("state=%s labels=%v", item.State, item.Labels)
	}
}
