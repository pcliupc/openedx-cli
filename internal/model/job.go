package model

// Job represents an asynchronous operation submitted to the backend.
type Job struct {
	JobID       string   `json:"job_id"`
	Operation   string   `json:"operation"`
	Status      string   `json:"status"`
	SubmittedAt string   `json:"submitted_at,omitempty"`
	FinishedAt  string   `json:"finished_at,omitempty"`
	Result      string   `json:"result,omitempty"`
	Artifacts   []string `json:"artifacts,omitempty"`
}
