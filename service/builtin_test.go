package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuiltinHtmlService(t *testing.T) {
	svc := BUILTIN_SERVICE_HTML
	f := tmpFile(t, "<h1>Test</h1>")
	out, err := svc.Run(f)
	assert.Nil(t, err)
	assert.Equal(t, "<h1>Test</h1>", string(out))
}

func TestBuiltinMarkdownService(t *testing.T) {
	svc := BUILTIN_SERVICE_MARKDOWN
	f := tmpFile(t, "# Test")
	out, err := svc.Run(f)
	assert.Nil(t, err)
	assert.Equal(t, "<h1 id=\"test\">Test</h1>\n", string(out))
}
