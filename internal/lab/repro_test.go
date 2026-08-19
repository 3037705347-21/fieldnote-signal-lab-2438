package lab

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListSiteFilterIsCaseInsensitive(t *testing.T) {
	service := NewService(NewMemoryStore(), DefaultCatalog())
	if _, err := service.Create(CreateObservationInput{
		Site:        "Wetland",
		ObservedAt:  time.Now().UTC().Add(-time.Hour).Format(time.RFC3339),
		Description: "Research team recorded repeated calls beside the marsh reeds.",
	}); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/api/v1/observations?site=wetland", nil)
	response := httptest.NewRecorder()
	NewServer(service).ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	var page Page
	if err := json.NewDecoder(response.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}
	if page.Total != 1 {
		t.Fatalf("page=%+v", page)
	}
}
