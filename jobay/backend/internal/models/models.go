package models

type JobStatus string

const (
	StatusDiscovered      JobStatus = "discovered"
	StatusQualified       JobStatus = "qualified"
	StatusReview          JobStatus = "review"
	StatusApplied         JobStatus = "applied"
	StatusOutcomeInterview JobStatus = "outcome_interview"
	StatusOutcomeRejected  JobStatus = "outcome_rejected"
)

type Job struct {
	ID           int       `json:"id"`
	Company      string    `json:"company"`
	Role         string    `json:"role"`
	URL          string    `json:"url,omitempty"`
	Status       JobStatus `json:"status"`
	Score        *int      `json:"score,omitempty"`
	Notes        string    `json:"notes,omitempty"`
	UserSlug     string    `json:"user_slug"`
	Location     string    `json:"location,omitempty"`
	Source       string    `json:"source,omitempty"`
	AppliedAt    string    `json:"applied_at,omitempty"`
	ApplyStatus  string    `json:"apply_status,omitempty"`
	MatchReasons string    `json:"match_reasons,omitempty"`
	CreatedAt    string    `json:"created_at"`
}

type User struct {
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	CVPath    string `json:"cv_path,omitempty"`
	CreatedAt string `json:"created_at"`
}

type Action struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Message   string `json:"message"`
	JobID     *int   `json:"job_id,omitempty"`
	UserSlug  string `json:"user_slug"`
	CreatedAt string `json:"created_at"`
}

type AgentStatus struct {
	ID         int    `json:"id"`
	Mode       string `json:"mode"`
	AIProvider string `json:"ai_provider"`
	IsRunning  int    `json:"is_running"`
	LastAction string `json:"last_action,omitempty"`
	UpdatedAt  string `json:"updated_at"`
}

type Run struct {
	ID             string `json:"id"`
	StartedAt      string `json:"started_at"`
	EndedAt        string `json:"ended_at,omitempty"`
	TotalFound     int    `json:"total_found"`
	TotalQualified int    `json:"total_qualified"`
	TotalApplied   int    `json:"total_applied"`
	TotalSkipped   int    `json:"total_skipped"`
}

type Stats struct {
	Discovered       int `json:"discovered"`
	Qualified        int `json:"qualified"`
	Review           int `json:"review"`
	Applied          int `json:"applied"`
	OutcomeInterview int `json:"outcome_interview"`
	OutcomeRejected  int `json:"outcome_rejected"`
	Total            int `json:"total"`
}

type Database struct {
	Jobs    []Job        `json:"jobs"`
	Users   []User       `json:"users"`
	Actions []Action     `json:"actions"`
	Agent   AgentStatus  `json:"agent"`
	Runs    []Run        `json:"runs"`
}
