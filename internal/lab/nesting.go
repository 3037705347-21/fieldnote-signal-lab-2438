package lab

import "time"

type NestingBand string

const (
	NestingBandLow      NestingBand = "low"
	NestingBandModerate NestingBand = "moderate"
	NestingBandHigh     NestingBand = "high"
)

type NestingEntry struct {
	ObservationID string      `json:"observation_id"`
	RecordedAt    time.Time   `json:"recorded_at"`
	Source        string      `json:"source"`
	Site          string      `json:"site"`
	Observer      string      `json:"observer"`
	Method        string      `json:"method"`
	Confidence    float64     `json:"confidence"`
	Band          NestingBand `json:"band"`
	Evidence      []string    `json:"evidence"`
	Notes         string      `json:"notes"`
	Revision      int         `json:"revision"`
	Archived      bool        `json:"archived"`
}
type NestingLedger struct {
	Entries     []NestingEntry `json:"entries"`
	Topic       string         `json:"topic"`
	Owner       string         `json:"owner"`
	Version     string         `json:"version"`
	Window      SamplingWindow `json:"window"`
	LastUpdated time.Time      `json:"last_updated"`
	Active      bool           `json:"active"`
	SourceCount int            `json:"source_count"`
	ReviewCount int            `json:"review_count"`
	Summary     string         `json:"summary"`
}
type NestingMetric struct {
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
type NestingProfile struct {
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
type NestingEvidence struct {
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
type NestingWindow struct {
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
type NestingReview struct {
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
type NestingSnapshot struct {
	ID            string      `json:"id"`
	Site          string      `json:"site"`
	CapturedAt    time.Time   `json:"captured_at"`
	Score         int         `json:"score"`
	Band          NestingBand `json:"band"`
	EvidenceCount int         `json:"evidence_count"`
	ReviewCount   int         `json:"review_count"`
	Summary       string      `json:"summary"`
	Stable        bool        `json:"stable"`
	Source        string      `json:"source"`
}
type NestingThreshold struct {
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
type NestingSummary struct {
	Site        string      `json:"site"`
	Topic       string      `json:"topic"`
	Total       int         `json:"total"`
	Reviewed    int         `json:"reviewed"`
	Score       int         `json:"score"`
	Band        NestingBand `json:"band"`
	WindowStart time.Time   `json:"window_start"`
	WindowEnd   time.Time   `json:"window_end"`
	Owner       string      `json:"owner"`
	GeneratedAt time.Time   `json:"generated_at"`
}

func NewNestingLedger() NestingLedger {
	return NestingLedger{Entries: []NestingEntry{}, Topic: "nesting", Active: true, Version: "v1", LastUpdated: nowUTC()}
}
func (l *NestingLedger) Add(entry NestingEntry) {
	l.Entries = append(l.Entries, entry)
	l.SourceCount = len(l.Entries)
	l.LastUpdated = nowUTC()
}
func (l NestingLedger) Count() int                     { return len(l.Entries) }
func (l NestingLedger) Ready() bool                    { return l.Active && len(l.Entries) > 0 }
func (l NestingLedger) OwnerName() string              { return l.Owner }
func (l NestingLedger) TopicName() string              { return l.Topic }
func (l NestingLedger) IsArchived() bool               { return !l.Active }
func (l NestingLedger) Updated() time.Time             { return l.LastUpdated }
func (e NestingEntry) Score() int                      { return int(e.Confidence * 100) }
func (e NestingEntry) IsReliable() bool                { return e.Confidence >= 0.5 && !e.Archived }
func (e NestingEntry) BandName() string                { return string(e.Band) }
func (e NestingEntry) EvidenceCount() int              { return len(e.Evidence) }
func (m NestingMetric) InRange() bool                  { return m.Value >= m.Minimum && m.Value <= m.Maximum }
func (m NestingMetric) IsRequired() bool               { return m.Required }
func (p NestingProfile) IsEnabled() bool               { return p.Enabled }
func (p NestingProfile) TagCount() int                 { return len(p.Tags) }
func (v NestingEvidence) IsReliable() bool             { return v.Reliable && v.Weight > 0 }
func (v NestingEvidence) Age() int                     { return int(nowUTC().Sub(v.CapturedAt).Hours()) }
func (w NestingWindow) IsComplete() bool               { return w.Complete && w.Observed >= w.Expected }
func (w NestingWindow) Duration() time.Duration        { return w.End.Sub(w.Start) }
func (r NestingReview) IsAccepted() bool               { return r.Accepted && r.Confidence >= 0.5 }
func (r NestingReview) NeedsFollowUp() bool            { return r.FollowUp }
func (s NestingSnapshot) IsStable() bool               { return s.Stable && s.Score >= 0 }
func (s NestingSnapshot) BandName() string             { return string(s.Band) }
func (t NestingThreshold) Contains(value float64) bool { return value >= t.Lower && value <= t.Upper }
func (t NestingThreshold) IsEnabled() bool             { return t.Enabled }
func (s NestingSummary) IsReviewed() bool              { return s.Reviewed <= s.Total }
func (s NestingSummary) WindowDuration() time.Duration { return s.WindowEnd.Sub(s.WindowStart) }
