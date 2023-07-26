package server

import (
	"io"
	"net/http"
	"testing"

	"github.com/apodhrad/tryit-editor/service"
	"github.com/stretchr/testify/assert"
)

func TestContentType(t *testing.T) {
	assert.Equal(t, CONTENT_TYPE_HTML, contentType("foo.htm"))
	assert.Equal(t, CONTENT_TYPE_HTML, contentType("foo.html"))
	assert.Equal(t, CONTENT_TYPE_CSS, contentType("foo.css"))
	assert.Equal(t, CONTENT_TYPE_JAVASCRIPT, contentType("foo.js"))
	assert.Equal(t, CONTENT_TYPE_TEXT, contentType("foo.txt"))
	assert.Equal(t, CONTENT_TYPE_TEXT, contentType("foo.text"))
	assert.Equal(t, CONTENT_TYPE_TEXT, contentType("foo.unknown"))

	// without extension
	assert.Equal(t, CONTENT_TYPE_TEXT, contentType("foo"))

	// with parent dir
	assert.Equal(t, CONTENT_TYPE_HTML, contentType("/dir/foo.html"))
}

const STATUS_200_OK string = "200 OK"
const STATUS_404_NOT_FOUND string = "404 Not Found"

func assertRequest(t *testing.T, addr string, path string, expected string) {
	url := "http://" + addr + path
	resp, err := http.Get(url)
	assert.Nil(t, err)
	if resp != nil {
		defer resp.Body.Close()
	}

	actual, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, expected, string(actual))
}

func assertResponseContain(t *testing.T, addr string, path string, expected string) {
	url := "http://" + addr + path
	resp, err := http.Get(url)
	assert.Nil(t, err)
	if resp != nil {
		defer resp.Body.Close()
	}

	actual, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Contains(t, string(actual), expected)
}

func assertStatus(t *testing.T, addr string, path string, expectedStatus string) {
	url := "http://" + addr + path
	resp, err := http.Get(url)
	if resp != nil {
		defer resp.Body.Close()
		assert.Equal(t, expectedStatus, resp.Status)
	} else {
		assert.Fail(t, err.Error())
	}
}

func TestServer(t *testing.T) {
	svc := service.SERVICE_CAT
	ctx, err := Start([]service.Service{svc})
	assert.Nil(t, err)
	assert.NotNil(t, ctx)

	assertStatus(t, "localhost:8080", "", STATUS_200_OK)
	assertStatus(t, "localhost:8080", "/", STATUS_200_OK)
	assertStatus(t, "localhost:8080", "/index.html", STATUS_200_OK)
	assertStatus(t, "localhost:8080", "/index.htmlx", STATUS_404_NOT_FOUND)

	assertStatus(t, "localhost:8080", "/test/test.txt", STATUS_200_OK)
	assertRequest(t, "localhost:8080", "/test/test.txt", "Hello World")

	assertResponseContain(t, "localhost:8080", "/index.html", "<option>"+svc.Name()+"</option>")

	err = Stop()
	assert.Nil(t, err)
}
