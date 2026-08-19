package lab

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListTotalIgnoresLimit(t *testing.T) {
	sequence = 0
	service := NewService(NewMemoryStore(), DefaultCatalog())
	for i := 0; i < 3; i++ {
		if _, err := service.Create(CreateObservationInput{
			Site:        "wetland",
			ObservedAt:  time.Now().UTC().Add(-time.Duration(i+1) * time.Hour).Format(time.RFC3339),
			Description: "Research team recorded repeated calls beside the marsh reeds.",
		}); err != nil {
			t.Fatal(err)
		}
	}

	request := httptest.NewRequest(http.MethodGet, "/api/v1/observations?site=wetland&limit=2", nil)
	response := httptest.NewRecorder()
	NewServer(service).ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	var page Page
	if err := json.NewDecoder(response.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}
	if page.Total != 3 || len(page.Items) != 2 {
		t.Fatalf("page=%+v", page)
	}
}
