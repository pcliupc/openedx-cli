package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openedx/cli/internal/auth"
	"github.com/openedx/cli/internal/config"
	"github.com/openedx/cli/internal/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubTokenProvider is a test double that always returns a fixed token.
type stubTokenProvider struct {
	token *auth.Token
	err   error
}

func (s *stubTokenProvider) Token(_ context.Context, _ config.Profile) (*auth.Token, error) {
	return s.token, s.err
}

// testProfile returns a profile configured to talk to the given baseURL.
func testProfile(baseURL string) config.Profile {
	return config.Profile{
		BaseURL:        baseURL,
		TokenURL:       baseURL + "/oauth2/token",
		ClientIDEnv:    "TEST_CLIENT_ID",
		ClientSecretEnv: "TEST_CLIENT_SECRET",
	}
}

// testCommand returns a sample command metadata for testing.
func testCommand() registry.CommandMeta {
	return registry.CommandMeta{
		Key:     "course.get",
		Method:  "GET",
		Path:    "/api/courses/v1/courses/{course_id}",
		OutputModel: "Course",
	}
}

// testPostCommand returns a sample POST command metadata.
func testPostCommand() registry.CommandMeta {
	return registry.CommandMeta{
		Key:          "course.create",
		Method:       "POST",
		Path:         "/api/courses/v1/courses",
		RequiredArgs: []string{"org", "number", "run", "title"},
		OutputModel:  "Course",
	}
}

func TestPublicProviderResolvesPath(t *testing.T) {
	result := resolvePath("/api/courses/v1/courses/{course_id}", map[string]string{
		"course_id": "course-v1:Org+Num+Run",
	})
	assert.Equal(t, "/api/courses/v1/courses/course-v1:Org+Num+Run", result)
}

func TestPublicProviderGETRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/api/courses/v1/courses/course-v1:Org+Num+Run", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"course-v1:Org+Num+Run","title":"Test Course"}`))
	}))
	defer server.Close()

	p := NewPublicProvider(server.Client(), testCommand())
	result, err := p.Execute(context.Background(), server.URL, "test-token", map[string]string{
		"course_id": "course-v1:Org+Num+Run",
	})

	require.NoError(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Contains(t, string(result.Body), "Test Course")
}

func TestPublicProviderPOSTRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/courses/v1/courses", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "TestOrg", body["org"])

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":"course-v1:TestOrg+Num+Run"}`))
	}))
	defer server.Close()

	p := NewPublicProvider(server.Client(), testPostCommand())
	result, err := p.Execute(context.Background(), server.URL, "test-token", map[string]string{
		"org": "TestOrg", "number": "Num", "run": "Run", "title": "Test",
	})

	require.NoError(t, err)
	assert.Equal(t, 201, result.StatusCode)
}

func TestPublicProviderReturnsProviderError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"detail":"Not found"}`))
	}))
	defer server.Close()

	p := NewPublicProvider(server.Client(), testCommand())
	_, err := p.Execute(context.Background(), server.URL, "test-token", map[string]string{
		"course_id": "nonexistent",
	})

	require.Error(t, err)
	provErr, ok := err.(*ProviderError)
	require.True(t, ok, "error should be ProviderError")
	assert.Equal(t, 404, provErr.StatusCode)
	assert.Contains(t, provErr.Message, "Not found")
}

func TestExtensionProviderGETRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/custom/api/courses", r.URL.Path)
		assert.Equal(t, "Bearer ext-token", r.Header.Get("Authorization"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"custom":true}`))
	}))
	defer server.Close()

	p := NewExtensionProvider(server.Client())
	result, err := p.Execute(context.Background(), "ext-token", config.ExtensionMapping{
		Method: "GET",
		URL:    server.URL + "/custom/api/courses",
	}, nil)

	require.NoError(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Contains(t, string(result.Body), "custom")
}

func TestExtensionProviderPOSTRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer ext-token", r.Header.Get("Authorization"))

		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "value1", body["key1"])

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	p := NewExtensionProvider(server.Client())
	result, err := p.Execute(context.Background(), "ext-token", config.ExtensionMapping{
		Method: "POST",
		URL:    server.URL + "/custom/api/courses",
	}, map[string]string{"key1": "value1"})

	require.NoError(t, err)
	assert.Equal(t, 200, result.StatusCode)
}

func TestExtensionProviderReturnsProviderError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal"}`))
	}))
	defer server.Close()

	p := NewExtensionProvider(server.Client())
	_, err := p.Execute(context.Background(), "ext-token", config.ExtensionMapping{
		Method: "GET",
		URL:    server.URL + "/custom/api",
	}, nil)

	require.Error(t, err)
	provErr, ok := err.(*ProviderError)
	require.True(t, ok)
	assert.Equal(t, 500, provErr.StatusCode)
}

func TestProviderErrorIsUnavailable(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		expected   bool
	}{
		{"404 Not Found", 404, true},
		{"405 Method Not Allowed", 405, true},
		{"501 Not Implemented", 501, true},
		{"400 Bad Request", 400, false},
		{"401 Unauthorized", 401, false},
		{"403 Forbidden", 403, false},
		{"500 Internal Server Error", 500, false},
		{"200 OK", 200, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := &ProviderError{StatusCode: tc.statusCode, Message: "test"}
			assert.Equal(t, tc.expected, err.IsUnavailable())
		})
	}
}

func TestProviderErrorImplementsError(t *testing.T) {
	err := &ProviderError{StatusCode: 404, Message: "not found"}
	assert.Equal(t, "not found", err.Error())
}

// --- Fallback tests ---

func TestFallbackUsesExtensionOnNotFound(t *testing.T) {
	// Public server returns 404.
	publicServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"detail":"Not found"}`))
	}))
	defer publicServer.Close()

	// Extension server returns 200.
	extServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"ext-course","source":"extension"}`))
	}))
	defer extServer.Close()

	fp := NewFallbackProvider(publicServer.Client())
	fp.Extension = NewExtensionProvider(extServer.Client())

	extMapping := &config.ExtensionMapping{
		Method: "GET",
		URL:    extServer.URL + "/ext/courses",
	}

	profile := testProfile(publicServer.URL)
	tp := &stubTokenProvider{token: &auth.Token{AccessToken: "test-token"}}

	result, err := fp.Execute(
		context.Background(), profile, tp,
		"course.get", testCommand(), extMapping,
		map[string]string{"course_id": "course-v1:Org+Num+Run"},
	)

	require.NoError(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Contains(t, string(result.Body), "extension")
}

func TestFallbackDoesNotUseExtensionOnForbidden(t *testing.T) {
	// Public server returns 403.
	publicServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"detail":"Forbidden"}`))
	}))
	defer publicServer.Close()

	extCallCount := 0
	extServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		extCallCount++
		w.WriteHeader(http.StatusOK)
	}))
	defer extServer.Close()

	fp := NewFallbackProvider(publicServer.Client())
	fp.Extension = NewExtensionProvider(extServer.Client())

	extMapping := &config.ExtensionMapping{
		Method: "GET",
		URL:    extServer.URL + "/ext/courses",
	}

	profile := testProfile(publicServer.URL)
	tp := &stubTokenProvider{token: &auth.Token{AccessToken: "test-token"}}

	_, err := fp.Execute(
		context.Background(), profile, tp,
		"course.get", testCommand(), extMapping,
		map[string]string{"course_id": "course-v1:Org+Num+Run"},
	)

	require.Error(t, err)
	provErr, ok := err.(*ProviderError)
	require.True(t, ok)
	assert.Equal(t, 403, provErr.StatusCode)
	assert.Equal(t, 0, extCallCount, "extension should NOT be called for 403")
}

func TestFallbackDoesNotUseExtensionOnValidationError(t *testing.T) {
	// Public server returns 400.
	publicServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"detail":"Bad request"}`))
	}))
	defer publicServer.Close()

	extCallCount := 0
	extServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		extCallCount++
		w.WriteHeader(http.StatusOK)
	}))
	defer extServer.Close()

	fp := NewFallbackProvider(publicServer.Client())
	fp.Extension = NewExtensionProvider(extServer.Client())

	extMapping := &config.ExtensionMapping{
		Method: "GET",
		URL:    extServer.URL + "/ext/courses",
	}

	profile := testProfile(publicServer.URL)
	tp := &stubTokenProvider{token: &auth.Token{AccessToken: "test-token"}}

	_, err := fp.Execute(
		context.Background(), profile, tp,
		"course.get", testCommand(), extMapping,
		map[string]string{"course_id": "invalid"},
	)

	require.Error(t, err)
	provErr, ok := err.(*ProviderError)
	require.True(t, ok)
	assert.Equal(t, 400, provErr.StatusCode)
	assert.Equal(t, 0, extCallCount, "extension should NOT be called for 400")
}

func TestFallbackDoesNotRetryWhenNoExtension(t *testing.T) {
	// Public server returns 404.
	publicServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"detail":"Not found"}`))
	}))
	defer publicServer.Close()

	fp := NewFallbackProvider(publicServer.Client())

	// No extension mapping (nil).
	profile := testProfile(publicServer.URL)
	tp := &stubTokenProvider{token: &auth.Token{AccessToken: "test-token"}}

	_, err := fp.Execute(
		context.Background(), profile, tp,
		"course.get", testCommand(), nil,
		map[string]string{"course_id": "course-v1:Org+Num+Run"},
	)

	require.Error(t, err)
	provErr, ok := err.(*ProviderError)
	require.True(t, ok)
	assert.Equal(t, 404, provErr.StatusCode)
}

func TestPublicSuccessReturnsDirectly(t *testing.T) {
	publicCallCount := 0
	// Public server returns 200.
	publicServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		publicCallCount++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"course-v1:Org+Num+Run","title":"Test"}`))
	}))
	defer publicServer.Close()

	extCallCount := 0
	extServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		extCallCount++
		w.WriteHeader(http.StatusOK)
	}))
	defer extServer.Close()

	fp := NewFallbackProvider(publicServer.Client())
	fp.Extension = NewExtensionProvider(extServer.Client())

	extMapping := &config.ExtensionMapping{
		Method: "GET",
		URL:    extServer.URL + "/ext/courses",
	}

	profile := testProfile(publicServer.URL)
	tp := &stubTokenProvider{token: &auth.Token{AccessToken: "test-token"}}

	result, err := fp.Execute(
		context.Background(), profile, tp,
		"course.get", testCommand(), extMapping,
		map[string]string{"course_id": "course-v1:Org+Num+Run"},
	)

	require.NoError(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Contains(t, string(result.Body), "Test")
	assert.Equal(t, 1, publicCallCount)
	assert.Equal(t, 0, extCallCount, "extension should NOT be called when public succeeds")
}

func TestFallbackUsesExtensionOnMethodNotAllowed(t *testing.T) {
	// Public server returns 405.
	publicServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"detail":"Method not allowed"}`))
	}))
	defer publicServer.Close()

	// Extension server returns 200.
	extServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"ext-course","source":"extension"}`))
	}))
	defer extServer.Close()

	fp := NewFallbackProvider(publicServer.Client())
	fp.Extension = NewExtensionProvider(extServer.Client())

	extMapping := &config.ExtensionMapping{
		Method: "GET",
		URL:    extServer.URL + "/ext/courses",
	}

	profile := testProfile(publicServer.URL)
	tp := &stubTokenProvider{token: &auth.Token{AccessToken: "test-token"}}

	result, err := fp.Execute(
		context.Background(), profile, tp,
		"course.get", testCommand(), extMapping,
		map[string]string{"course_id": "course-v1:Org+Num+Run"},
	)

	require.NoError(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Contains(t, string(result.Body), "extension")
}

func TestFallbackUsesExtensionOnNotImplemented(t *testing.T) {
	// Public server returns 501.
	publicServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(`{"detail":"Not implemented"}`))
	}))
	defer publicServer.Close()

	// Extension server returns 200.
	extServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"ext-course","source":"extension"}`))
	}))
	defer extServer.Close()

	fp := NewFallbackProvider(publicServer.Client())
	fp.Extension = NewExtensionProvider(extServer.Client())

	extMapping := &config.ExtensionMapping{
		Method: "GET",
		URL:    extServer.URL + "/ext/courses",
	}

	profile := testProfile(publicServer.URL)
	tp := &stubTokenProvider{token: &auth.Token{AccessToken: "test-token"}}

	result, err := fp.Execute(
		context.Background(), profile, tp,
		"course.get", testCommand(), extMapping,
		map[string]string{"course_id": "course-v1:Org+Num+Run"},
	)

	require.NoError(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Contains(t, string(result.Body), "extension")
}
