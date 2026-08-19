package lab

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListTagFilterTrimsWhitespace(t *testing.T) {
	service := NewService(NewMemoryStore(), DefaultCatalog())
	SeedDemo(service)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/observations?tag=%20acoustic%20", nil)
	response := httptest.NewRecorder()
	NewServer(service).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	var page Page
	if err := json.NewDecoder(response.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}
	if page.Total != 1 || len(page.Items) != 1 || !hasLabel(page.Items[0], "acoustic") {
		t.Fatalf("page=%+v", page)
	}
}
