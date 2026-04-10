package normalize

import (
	"os"
	"testing"

	"github.com/openedx/cli/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCertificateListFromPublicPayload(t *testing.T) {
	data, err := os.ReadFile("../../testdata/public/certificate_list.json")
	require.NoError(t, err)

	certs, err := CertificateListFromJSON(data)
	require.NoError(t, err)

	require.Len(t, certs, 1)
	expected := &model.Certificate{
		Username:        "alice",
		CourseID:        "course-v1:demo+cs101+2026",
		CertificateType: "verified",
		Status:          "downloadable",
		DownloadURL:     "https://openedx.example.com/certificates/abc123",
		Grade:           "85%",
	}
	assert.Equal(t, expected, certs[0])
}

func TestCertificateListFromExtensionPayload(t *testing.T) {
	data, err := os.ReadFile("../../testdata/extension/certificate_list.json")
	require.NoError(t, err)

	certs, err := CertificateListFromJSON(data)
	require.NoError(t, err)

	require.Len(t, certs, 1)
	expected := &model.Certificate{
		Username:        "alice",
		CourseID:        "course-v1:demo+cs101+2026",
		CertificateType: "verified",
		Status:          "downloadable",
	}
	assert.Equal(t, expected, certs[0])
}

func TestCertificateListFromJSONInvalidInput(t *testing.T) {
	_, err := CertificateListFromJSON([]byte("not json"))
	assert.Error(t, err)
}
