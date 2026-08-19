package lab

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListStateFilterIsCaseInsensitive(t *testing.T) {
	service := NewService(NewMemoryStore(), DefaultCatalog())
	SeedDemo(service)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/observations?state=Reviewed", nil)
	response := httptest.NewRecorder()
	NewServer(service).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	var page Page
	if err := json.NewDecoder(response.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}
	if page.Total != 1 || len(page.Items) != 1 || page.Items[0].State != StateReviewed {
		t.Fatalf("page=%+v", page)
	}
}
