package lab

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListNormalizesTagQuery(t *testing.T) {
	sequence = 0
	service := NewService(NewMemoryStore(), DefaultCatalog())
	item, err := service.Create(CreateObservationInput{Site: "wetland", ObservedAt: time.Now().UTC().Add(-time.Hour).Format(time.RFC3339), Description: "Research team logged repeated calls beside the marsh reeds before sunrise."})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = service.AddLabels(item.ID, AddLabelsInput{Labels: []string{"acoustic"}}); err != nil {
		t.Fatal(err)
	}
	handler := NewServer(service)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/observations?tag=Acoustic", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var page Page
	if err := json.NewDecoder(rr.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}
	if page.Total != 1 || len(page.Items) != 1 || page.Items[0].ID != item.ID {
		t.Fatalf("page=%+v", page)
	}
}
