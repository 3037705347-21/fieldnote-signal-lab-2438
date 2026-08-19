package lab

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReportIncludesFullLabelContribution(t *testing.T) {
	service := NewService(NewMemoryStore(), DefaultCatalog())
	SeedDemo(service)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/observations/demo-001/report", nil)
	response := httptest.NewRecorder()
	NewServer(service).ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	var report SignalReport
	if err := json.NewDecoder(response.Body).Decode(&report); err != nil {
		t.Fatal(err)
	}
	if report.Score != 79 {
		t.Fatalf("report=%+v", report)
	}
}
