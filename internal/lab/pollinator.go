package lab

import "time"

type PollinatorBand string

const (
	PollinatorBandLow      PollinatorBand = "low"
	PollinatorBandModerate PollinatorBand = "moderate"
	PollinatorBandHigh     PollinatorBand = "high"
)

type PollinatorEntry struct {
	ObservationID string         `json:"observation_id"`
	RecordedAt    time.Time      `json:"recorded_at"`
	Source        string         `json:"source"`
	Site          string         `json:"site"`
	Observer      string         `json:"observer"`
	Method        string         `json:"method"`
	Confidence    float64        `json:"confidence"`
	Band          PollinatorBand `json:"band"`
	Evidence      []string       `json:"evidence"`
	Notes         string         `json:"notes"`
	Revision      int            `json:"revision"`
	Archived      bool           `json:"archived"`
}
type PollinatorLedger struct {
	Entries     []PollinatorEntry `json:"entries"`
	Topic       string            `json:"topic"`
	Owner       string            `json:"owner"`
	Version     string            `json:"version"`
	Window      SamplingWindow    `json:"window"`
	LastUpdated time.Time         `json:"last_updated"`
	Active      bool              `json:"active"`
	SourceCount int               `json:"source_count"`
	ReviewCount int               `json:"review_count"`
	Summary     string            `json:"summary"`
}
type PollinatorMetric struct {
	Name        string    `json:"name"`
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Weight      int       `json:"weight"`
	Required    bool      `json:"required"`
	Explanation string    `json:"explanation"`
	Minimum     float64   `json:"minimum"`
	Maximum     float64   `json:"maximum"`
	UpdatedAt   time.Time `json:"updated_at"`
	Topic       string    `json:"topic"`
}
type PollinatorProfile struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Site        string    `json:"site"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Enabled     bool      `json:"enabled"`
	Priority    int       `json:"priority"`
	Tags        []string  `json:"tags"`
	Owner       string    `json:"owner"`
}
type PollinatorEvidence struct {
	ID         string    `json:"id"`
	EntryID    string    `json:"entry_id"`
	Kind       string    `json:"kind"`
	Value      string    `json:"value"`
	CapturedAt time.Time `json:"captured_at"`
	Reliable   bool      `json:"reliable"`
	Weight     int       `json:"weight"`
	Note       string    `json:"note"`
	Author     string    `json:"author"`
	Revision   int       `json:"revision"`
}
type PollinatorWindow struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Site     string    `json:"site"`
	Label    string    `json:"label"`
	Expected int       `json:"expected"`
	Observed int       `json:"observed"`
	Complete bool      `json:"complete"`
	Closed   bool      `json:"closed"`
	Owner    string    `json:"owner"`
	Topic    string    `json:"topic"`
}
type PollinatorReview struct {
	EntryID    string    `json:"entry_id"`
	Verdict    string    `json:"verdict"`
	Confidence float64   `json:"confidence"`
	Reviewer   string    `json:"reviewer"`
	ReviewedAt time.Time `json:"reviewed_at"`
	Notes      string    `json:"notes"`
	Accepted   bool      `json:"accepted"`
	FollowUp   bool      `json:"follow_up"`
	Queue      string    `json:"queue"`
	Version    string    `json:"version"`
}
type PollinatorSnapshot struct {
	ID            string         `json:"id"`
	Site          string         `json:"site"`
	CapturedAt    time.Time      `json:"captured_at"`
	Score         int            `json:"score"`
	Band          PollinatorBand `json:"band"`
	EvidenceCount int            `json:"evidence_count"`
	ReviewCount   int            `json:"review_count"`
	Summary       string         `json:"summary"`
	Stable        bool           `json:"stable"`
	Source        string         `json:"source"`
}
type PollinatorThreshold struct {
	Name        string    `json:"name"`
	Lower       float64   `json:"lower"`
	Upper       float64   `json:"upper"`
	Unit        string    `json:"unit"`
	Strict      bool      `json:"strict"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	Owner       string    `json:"owner"`
	UpdatedAt   time.Time `json:"updated_at"`
	Revision    int       `json:"revision"`
}
type PollinatorSummary struct {
	Site        string         `json:"site"`
	Topic       string         `json:"topic"`
	Total       int            `json:"total"`
	Reviewed    int            `json:"reviewed"`
	Score       int            `json:"score"`
	Band        PollinatorBand `json:"band"`
	WindowStart time.Time      `json:"window_start"`
	WindowEnd   time.Time      `json:"window_end"`
	Owner       string         `json:"owner"`
	GeneratedAt time.Time      `json:"generated_at"`
}

func NewPollinatorLedger() PollinatorLedger {
	return PollinatorLedger{Entries: []PollinatorEntry{}, Topic: "pollinator", Active: true, Version: "v1", LastUpdated: nowUTC()}
}
func (l *PollinatorLedger) Add(entry PollinatorEntry) {
	l.Entries = append(l.Entries, entry)
	l.SourceCount = len(l.Entries)
	l.LastUpdated = nowUTC()
}
func (l PollinatorLedger) Count() int              { return len(l.Entries) }
func (l PollinatorLedger) Ready() bool             { return l.Active && len(l.Entries) > 0 }
func (l PollinatorLedger) OwnerName() string       { return l.Owner }
func (l PollinatorLedger) TopicName() string       { return l.Topic }
func (l PollinatorLedger) IsArchived() bool        { return !l.Active }
func (l PollinatorLedger) Updated() time.Time      { return l.LastUpdated }
func (e PollinatorEntry) Score() int               { return int(e.Confidence * 100) }
func (e PollinatorEntry) IsReliable() bool         { return e.Confidence >= 0.5 && !e.Archived }
func (e PollinatorEntry) BandName() string         { return string(e.Band) }
func (e PollinatorEntry) EvidenceCount() int       { return len(e.Evidence) }
func (m PollinatorMetric) InRange() bool           { return m.Value >= m.Minimum && m.Value <= m.Maximum }
func (m PollinatorMetric) IsRequired() bool        { return m.Required }
func (p PollinatorProfile) IsEnabled() bool        { return p.Enabled }
func (p PollinatorProfile) TagCount() int          { return len(p.Tags) }
func (v PollinatorEvidence) IsReliable() bool      { return v.Reliable && v.Weight > 0 }
func (v PollinatorEvidence) Age() int              { return int(nowUTC().Sub(v.CapturedAt).Hours()) }
func (w PollinatorWindow) IsComplete() bool        { return w.Complete && w.Observed >= w.Expected }
func (w PollinatorWindow) Duration() time.Duration { return w.End.Sub(w.Start) }
func (r PollinatorReview) IsAccepted() bool        { return r.Accepted && r.Confidence >= 0.5 }
func (r PollinatorReview) NeedsFollowUp() bool     { return r.FollowUp }
func (s PollinatorSnapshot) IsStable() bool        { return s.Stable && s.Score >= 0 }
func (s PollinatorSnapshot) BandName() string      { return string(s.Band) }
func (t PollinatorThreshold) Contains(value float64) bool {
	return value >= t.Lower && value <= t.Upper
}
func (t PollinatorThreshold) IsEnabled() bool             { return t.Enabled }
func (s PollinatorSummary) IsReviewed() bool              { return s.Reviewed <= s.Total }
func (s PollinatorSummary) WindowDuration() time.Duration { return s.WindowEnd.Sub(s.WindowStart) }
