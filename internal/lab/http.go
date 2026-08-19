package lab

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	service *Service
	mux     *http.ServeMux
}

func NewServer(service *Service) http.Handler {
	server := &Server{service: service, mux: http.NewServeMux()}
	server.routes()
	return server.withLog(server.mux)
}
func (s *Server) routes() {
	s.mux.HandleFunc("GET /health", s.health)
	s.mux.HandleFunc("GET /api/v1/catalog", s.catalog)
	s.mux.HandleFunc("GET /api/v1/observations", s.list)
	s.mux.HandleFunc("POST /api/v1/observations", s.create)
	s.mux.HandleFunc("GET /api/v1/observations/{id}", s.get)
	s.mux.HandleFunc("POST /api/v1/observations/{id}/labels", s.labels)
	s.mux.HandleFunc("POST /api/v1/observations/{id}/review", s.review)
	s.mux.HandleFunc("GET /api/v1/observations/{id}/report", s.report)
}
func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]string{"status": "ok", "service": "fieldnote-signal-lab"})
}
func (s *Server) catalog(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, s.service.Catalog())
}
func (s *Server) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	page, err := s.service.List(q.Get("site"), q.Get("state"), strings.ToLower(q.Get("tag")), limit)
	if err != nil {
		writeFailure(w, 400, err)
		return
	}
	writeJSON(w, 200, page)
}
func (s *Server) create(w http.ResponseWriter, r *http.Request) {
	var input CreateObservationInput
	if err := decode(r, &input); err != nil {
		writeFailure(w, 400, err)
		return
	}
	item, err := s.service.Create(input)
	if err != nil {
		writeFailure(w, 400, err)
		return
	}
	writeJSON(w, 201, item)
}
func (s *Server) get(w http.ResponseWriter, r *http.Request) {
	item, err := s.service.Get(r.PathValue("id"))
	if err != nil {
		writeFailure(w, 404, err)
		return
	}
	writeJSON(w, 200, item)
}
func (s *Server) labels(w http.ResponseWriter, r *http.Request) {
	var input AddLabelsInput
	if err := decode(r, &input); err != nil {
		writeFailure(w, 400, err)
		return
	}
	item, err := s.service.AddLabels(r.PathValue("id"), input)
	if err != nil {
		writeFailure(w, statusFor(err), err)
		return
	}
	writeJSON(w, 200, item)
}
func (s *Server) review(w http.ResponseWriter, r *http.Request) {
	var input ReviewInput
	if err := decode(r, &input); err != nil {
		writeFailure(w, 400, err)
		return
	}
	item, err := s.service.Review(r.PathValue("id"), input)
	if err != nil {
		writeFailure(w, statusFor(err), err)
		return
	}
	writeJSON(w, 200, item)
}
func (s *Server) report(w http.ResponseWriter, r *http.Request) {
	report, err := s.service.Report(r.PathValue("id"))
	if err != nil {
		writeFailure(w, 404, err)
		return
	}
	writeJSON(w, 200, report)
}
func decode(r *http.Request, target any) error {
	return decodeStrictJSON(r.Body, target)
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeFailure(w http.ResponseWriter, status int, err error) {
	code := "invalid_request"
	if errors.Is(err, ErrNotFound) {
		code = "not_found"
	}
	if errors.Is(err, ErrConflict) {
		code = "conflict"
	}
	writeJSON(w, status, APIError{Code: code, Message: err.Error()})
}
func statusFor(err error) int {
	if errors.Is(err, ErrNotFound) {
		return 404
	}
	return 400
}
func (s *Server) withLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { started := time.Now(); next.ServeHTTP(w, r); _ = started })
}
