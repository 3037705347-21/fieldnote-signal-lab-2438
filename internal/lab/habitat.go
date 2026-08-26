package lab

import "time"

type HabitatBand string

const (
	HabitatBandLow      HabitatBand = "low"
	HabitatBandModerate HabitatBand = "moderate"
	HabitatBandHigh     HabitatBand = "high"
)

type HabitatEntry struct {
	ObservationID string      `json:"observation_id"`
	RecordedAt    time.Time   `json:"recorded_at"`
	Source        string      `json:"source"`
	Site          string      `json:"site"`
	Observer      string      `json:"observer"`
	Method        string      `json:"method"`
	Confidence    float64     `json:"confidence"`
	Band          HabitatBand `json:"band"`
	Evidence      []string    `json:"evidence"`
	Notes         string      `json:"notes"`
	Revision      int         `json:"revision"`
	Archived      bool        `json:"archived"`
}
type HabitatLedger struct {
	Entries     []HabitatEntry `json:"entries"`
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
type HabitatMetric struct {
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
type HabitatProfile struct {
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
type HabitatEvidence struct {
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
type HabitatWindow struct {
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
type HabitatReview struct {
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
type HabitatSnapshot struct {
	ID            string      `json:"id"`
	Site          string      `json:"site"`
	CapturedAt    time.Time   `json:"captured_at"`
	Score         int         `json:"score"`
	Band          HabitatBand `json:"band"`
	EvidenceCount int         `json:"evidence_count"`
	ReviewCount   int         `json:"review_count"`
	Summary       string      `json:"summary"`
	Stable        bool        `json:"stable"`
	Source        string      `json:"source"`
}
type HabitatThreshold struct {
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
type HabitatSummary struct {
	Site        string      `json:"site"`
	Topic       string      `json:"topic"`
	Total       int         `json:"total"`
	Reviewed    int         `json:"reviewed"`
	Score       int         `json:"score"`
	Band        HabitatBand `json:"band"`
	WindowStart time.Time   `json:"window_start"`
	WindowEnd   time.Time   `json:"window_end"`
	Owner       string      `json:"owner"`
	GeneratedAt time.Time   `json:"generated_at"`
}

func NewHabitatLedger() HabitatLedger {
	return HabitatLedger{Entries: []HabitatEntry{}, Topic: "habitat", Active: true, Version: "v1", LastUpdated: nowUTC()}
}
func (l *HabitatLedger) Add(entry HabitatEntry) {
	l.Entries = append(l.Entries, entry)
	l.SourceCount = len(l.Entries)
	l.LastUpdated = nowUTC()
}
func (l HabitatLedger) Count() int                     { return len(l.Entries) }
func (l HabitatLedger) Ready() bool                    { return l.Active && len(l.Entries) > 0 }
func (l HabitatLedger) OwnerName() string              { return l.Owner }
func (l HabitatLedger) TopicName() string              { return l.Topic }
func (l HabitatLedger) IsArchived() bool               { return !l.Active }
func (l HabitatLedger) Updated() time.Time             { return l.LastUpdated }
func (e HabitatEntry) Score() int                      { return int(e.Confidence * 100) }
func (e HabitatEntry) IsReliable() bool                { return e.Confidence >= 0.5 && !e.Archived }
func (e HabitatEntry) BandName() string                { return string(e.Band) }
func (e HabitatEntry) EvidenceCount() int              { return len(e.Evidence) }
func (m HabitatMetric) InRange() bool                  { return m.Value >= m.Minimum && m.Value <= m.Maximum }
func (m HabitatMetric) IsRequired() bool               { return m.Required }
func (p HabitatProfile) IsEnabled() bool               { return p.Enabled }
func (p HabitatProfile) TagCount() int                 { return len(p.Tags) }
func (v HabitatEvidence) IsReliable() bool             { return v.Reliable && v.Weight > 0 }
func (v HabitatEvidence) Age() int                     { return int(nowUTC().Sub(v.CapturedAt).Hours()) }
func (w HabitatWindow) IsComplete() bool               { return w.Complete && w.Observed >= w.Expected }
func (w HabitatWindow) Duration() time.Duration        { return w.End.Sub(w.Start) }
func (r HabitatReview) IsAccepted() bool               { return r.Accepted && r.Confidence >= 0.5 }
func (r HabitatReview) NeedsFollowUp() bool            { return r.FollowUp }
func (s HabitatSnapshot) IsStable() bool               { return s.Stable && s.Score >= 0 }
func (s HabitatSnapshot) BandName() string             { return string(s.Band) }
func (t HabitatThreshold) Contains(value float64) bool { return value >= t.Lower && value <= t.Upper }
func (t HabitatThreshold) IsEnabled() bool             { return t.Enabled }
func (s HabitatSummary) IsReviewed() bool              { return s.Reviewed <= s.Total }
func (s HabitatSummary) WindowDuration() time.Duration { return s.WindowEnd.Sub(s.WindowStart) }
