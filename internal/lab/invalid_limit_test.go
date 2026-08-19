package lab

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRejectsInvalidListLimit(t *testing.T) {
	service := NewService(NewMemoryStore(), DefaultCatalog())
	SeedDemo(service)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/observations?limit=not-a-number", nil)
	recorder := httptest.NewRecorder()

	NewServer(service).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status=%d, want %d", recorder.Code, http.StatusBadRequest)
	}
}
