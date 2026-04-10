package model

// Certificate represents a user's course completion certificate.
type Certificate struct {
	Username       string `json:"username"`
	CourseID       string `json:"course_id"`
	CertificateType string `json:"certificate_type"`
	Status         string `json:"status"`
	DownloadURL    string `json:"download_url,omitempty"`
	Grade          string `json:"grade,omitempty"`
}
