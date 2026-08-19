package lab

import (
	"strings"
	"time"
)

type ObservationState string

const (
	StateCaptured ObservationState = "captured"
	StateLabeled  ObservationState = "labeled"
	StateReviewed ObservationState = "reviewed"
)

type Observation struct {
	ID          string           `json:"id"`
	Site        string           `json:"site"`
	ObservedAt  time.Time        `json:"observed_at"`
	Description string           `json:"description"`
	State       ObservationState `json:"state"`
	Labels      []string         `json:"labels"`
	Review      *Review          `json:"review,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}
type CreateObservationInput struct {
	Site        string `json:"site"`
	ObservedAt  string `json:"observed_at"`
	Description string `json:"description"`
}
type AddLabelsInput struct {
	Labels []string `json:"labels"`
}
type ReviewInput struct {
	Verdict    string  `json:"verdict"`
	Confidence float64 `json:"confidence"`
	Notes      string  `json:"notes"`
}
type Review struct {
	Verdict    string    `json:"verdict"`
	Confidence float64   `json:"confidence"`
	Notes      string    `json:"notes"`
	ReviewedAt time.Time `json:"reviewed_at"`
}
type SignalReport struct {
	ObservationID string    `json:"observation_id"`
	Score         int       `json:"score"`
	Band          string    `json:"band"`
	Reasons       []string  `json:"reasons"`
	GeneratedAt   time.Time `json:"generated_at"`
}
type Catalog struct {
	Habitats []string `json:"habitats"`
	Tags     []string `json:"tags"`
	Verdicts []string `json:"verdicts"`
}
type SamplingWindow struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Site  string    `json:"site"`
}
type Page struct {
	Items []Observation `json:"items"`
	Total int           `json:"total"`
}
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func normalizeFilterValue(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
