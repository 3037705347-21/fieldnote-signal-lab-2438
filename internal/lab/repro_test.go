package lab

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListAppliesLimitAfterFiltering(t *testing.T) {
	sequence = 0
	service := NewService(NewMemoryStore(), DefaultCatalog())
	for _, site := range []string{"wetland", "forest"} {
		if _, err := service.Create(CreateObservationInput{Site: site, ObservedAt: time.Now().UTC().Add(-time.Hour).Format(time.RFC3339), Description: "Research team recorded a sustained observation along the marked transect."}); err != nil {
			t.Fatal(err)
		}
	}
	handler := NewServer(service)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/observations?site=forest&limit=1", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var page Page
	if err := json.NewDecoder(rr.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}
	if page.Total != 1 || len(page.Items) != 1 || page.Items[0].Site != "forest" {
		t.Fatalf("page=%+v", page)
	}
}
