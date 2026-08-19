package lab

import (
	"encoding/json"
	"errors"
	"io"
	"math"
	"sort"
	"strings"
	"time"
)

var ErrNotFound = errors.New("observation not found")
var ErrConflict = errors.New("observation already exists")

type FieldError struct {
	Field   string
	Message string
}

func (e FieldError) Error() string        { return e.Field + ": " + e.Message }
func invalid(field, message string) error { return FieldError{field, message} }
func validateCreate(input CreateObservationInput) (time.Time, error) {
	if strings.TrimSpace(input.Site) == "" {
		return time.Time{}, invalid("site", "is required")
	}
	if len(strings.TrimSpace(input.Description)) < 12 {
		return time.Time{}, invalid("description", "must contain at least 12 characters")
	}
	value, err := time.Parse(time.RFC3339, input.ObservedAt)
	if err != nil {
		return time.Time{}, invalid("observed_at", "must be RFC3339")
	}
	if value.After(time.Now().UTC().Add(2 * time.Minute)) {
		return time.Time{}, invalid("observed_at", "cannot be in the future")
	}
	return value.UTC(), nil
}
func validateLabels(c Catalog, values []string) ([]string, error) {
	labels := normalizeLabels(values)
	if len(labels) == 0 {
		return nil, invalid("labels", "must contain at least one tag")
	}
	for _, label := range labels {
		if !contains(c.Tags, label) {
			return nil, invalid("labels", "contains unsupported tag "+label)
		}
	}
	return labels, nil
}
func validateReview(c Catalog, input ReviewInput) error {
	if !contains(c.Verdicts, strings.TrimSpace(input.Verdict)) {
		return invalid("verdict", "is unsupported")
	}
	if math.IsNaN(input.Confidence) || math.IsInf(input.Confidence, 0) || input.Confidence < 0 || input.Confidence > 1 {
		return invalid("confidence", "must be between 0 and 1")
	}
	if len(strings.TrimSpace(input.Notes)) < 8 {
		return invalid("notes", "must contain at least 8 characters")
	}
	return nil
}
func normalizeLabels(values []string) []string {
	seen := map[string]bool{}
	result := []string{}
	for _, value := range values {
		label := strings.ToLower(strings.TrimSpace(value))
		if label != "" && !seen[label] {
			seen[label] = true
			result = append(result, label)
		}
	}
	sort.Strings(result)
	return result
}
func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}
func validState(state ObservationState) bool {
	return state == StateCaptured || state == StateLabeled || state == StateReviewed
}
func normalizeCreateInput(input CreateObservationInput) CreateObservationInput {
	input.Site = strings.TrimSpace(input.Site)
	input.Description = strings.TrimSpace(input.Description)
	return input
}
func decodeStrictJSON(body io.Reader, target any) error {
	decoder := json.NewDecoder(io.LimitReader(body, maxJSONBodyBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return errors.New("request body must contain one JSON value")
		}
		return err
	}
	return nil
}
