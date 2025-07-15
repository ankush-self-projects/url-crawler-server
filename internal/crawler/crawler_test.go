package crawler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectHTMLVersion(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "HTML5",
			html:     "<!doctype html><html><head><title>Test</title></head><body></body></html>",
			expected: "HTML5",
		},
		{
			name:     "Older HTML",
			html:     "<html><head><title>Test</title></head><body></body></html>",
			expected: "Older HTML",
		},
		{
			name:     "HTML4",
			html:     "<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01//EN\"><html><head><title>Test</title></head><body></body></html>",
			expected: "Older HTML",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectHTMLVersion(tt.html)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractTitle(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "Simple title",
			html:     "<html><head><title>Test Page</title></head><body></body></html>",
			expected: "Test Page",
		},
		{
			name:     "No title",
			html:     "<html><head></head><body></body></html>",
			expected: "",
		},
		{
			name:     "Multiple titles (should get first)",
			html:     "<html><head><title>First Title</title><title>Second Title</title></head><body></body></html>",
			expected: "First Title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parseDocument(tt.html)
			assert.NoError(t, err)
			result := extractTitle(doc)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractHeadingSummary(t *testing.T) {
	html := `
		<html>
			<head><title>Test</title></head>
			<body>
				<h1>Main Title</h1>
				<h2>Subtitle 1</h2>
				<h2>Subtitle 2</h2>
				<h3>Sub-subtitle</h3>
				<p>Some text</p>
			</body>
		</html>
	`

	doc, err := parseDocument(html)
	assert.NoError(t, err)

	result := extractHeadingSummary(doc)
	expected := "H1: 1, H2: 2, H3: 1"
	assert.Equal(t, expected, result)
}

func TestDetectLoginForm(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected bool
	}{
		{
			name: "Has login form",
			html: `
				<html>
					<head><title>Login</title></head>
					<body>
						<form>
							<input type="text" name="username">
							<input type="password" name="password">
							<button type="submit">Login</button>
						</form>
					</body>
				</html>
			`,
			expected: true,
		},
		{
			name: "No login form",
			html: `
				<html>
					<head><title>Home</title></head>
					<body>
						<form>
							<input type="text" name="search">
							<button type="submit">Search</button>
						</form>
					</body>
				</html>
			`,
			expected: false,
		},
		{
			name: "No forms",
			html: `
				<html>
					<head><title>Home</title></head>
					<body>
						<h1>Welcome</h1>
						<p>Some content</p>
					</body>
				</html>
			`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parseDocument(tt.html)
			assert.NoError(t, err)
			result := detectLoginForm(doc)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAnalyzeLinks(t *testing.T) {
	html := `
		<html>
			<head><title>Test</title></head>
			<body>
				<a href="https://example.com/page1">Internal Link 1</a>
				<a href="https://example.com/page2">Internal Link 2</a>
				<a href="https://google.com">External Link 1</a>
				<a href="https://github.com">External Link 2</a>
				<a href="https://broken-link-that-doesnt-exist.com">Broken Link</a>
			</body>
		</html>
	`

	doc, err := parseDocument(html)
	assert.NoError(t, err)

	internal, external, broken := analyzeLinks(doc, "https://example.com")

	// Note: The broken link count might vary depending on network conditions
	// We'll just check that we have the expected internal and external counts
	assert.Equal(t, 2, internal)
	assert.Equal(t, 2, external)
	// broken count could be 0 or 1 depending on network
	assert.GreaterOrEqual(t, broken, 0)
	assert.LessOrEqual(t, broken, 1)
}
