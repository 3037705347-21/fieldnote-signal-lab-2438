package lab

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRejectsTrailingJSON(t *testing.T) {
	service := NewService(NewMemoryStore(), DefaultCatalog())
	body := `{"site":"wetland","observed_at":"2026-08-18T10:00:00Z","description":"Research team recorded repeated calls beside the marsh reeds."}{"extra":true}`
	request := httptest.NewRequest(http.MethodPost, "/api/v1/observations", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	NewServer(service).ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
}
