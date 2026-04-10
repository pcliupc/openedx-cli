package normalize

import (
	"encoding/json"

	"github.com/openedx/cli/internal/model"
)

// rawCertificate is a lax intermediate struct that accepts field names from both
// public and extension API payloads.
type rawCertificate struct {
	Username        string `json:"username"`
	CourseID        string `json:"course_id"`
	CourseKey       string `json:"course_key"`
	CertificateType string `json:"certificate_type"`
	CertType        string `json:"cert_type"`
	Type            string `json:"type"`
	Status          string `json:"status"`
	DownloadURL     string `json:"download_url"`
	Grade           string `json:"grade"`
}

// rawCertificateList wraps the public API paginated response shape.
type rawCertificateList struct {
	Results []rawCertificate `json:"results"`
}

func (r rawCertificate) toModel() *model.Certificate {
	courseID := firstNonEmpty(r.CourseID, r.CourseKey)
	certType := firstNonEmpty(r.CertificateType, r.CertType, r.Type)
	return &model.Certificate{
		Username:        r.Username,
		CourseID:        courseID,
		CertificateType: certType,
		Status:          r.Status,
		DownloadURL:     r.DownloadURL,
		Grade:           r.Grade,
	}
}

// CertificateListFromJSON parses a list of certificates from raw JSON bytes.
// It supports the public API paginated response {"results": [...]} and the
// extension API top-level array [...].
func CertificateListFromJSON(data []byte) ([]*model.Certificate, error) {
	// Try paginated public API shape first.
	var paginated rawCertificateList
	if err := json.Unmarshal(data, &paginated); err == nil && len(paginated.Results) > 0 {
		out := make([]*model.Certificate, len(paginated.Results))
		for i, r := range paginated.Results {
			out[i] = r.toModel()
		}
		return out, nil
	}

	// Try top-level array (extension API shape).
	var items []rawCertificate
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	out := make([]*model.Certificate, len(items))
	for i, r := range items {
		out[i] = r.toModel()
	}
	return out, nil
}
