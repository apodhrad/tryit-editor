package server

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

const DEFAULT_PORT int = 8080

const CONTENT_TYPE_CSS string = "text/css"
const CONTENT_TYPE_TEXT string = "text/plain"
const CONTENT_TYPE_HTML string = "text/html"
const CONTENT_TYPE_JAVASCRIPT string = "text/javascript"
const CONTENT_TYPE_JSON string = "application/json"
const CONTENT_TYPE_FILE string = "application/octet-stream"

//go:embed html
var htmlFS embed.FS

func contentType(path string) string {
	ext := filepath.Ext(path)
	switch ext {
	case ".htm":
		return CONTENT_TYPE_HTML
	case ".html":
		return CONTENT_TYPE_HTML
	case ".js":
		return CONTENT_TYPE_JAVASCRIPT
	case ".css":
		return CONTENT_TYPE_CSS
	default:
		return CONTENT_TYPE_TEXT
	}
}

func htmlHandler(fsys fs.FS, path string) (string, func(w http.ResponseWriter, r *http.Request)) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			fmt.Println("[ERROR] " + err.Error())
			w.Header().Set("Content-Type", CONTENT_TYPE_TEXT)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
		} else {
			w.Header().Set("Content-Type", contentType(path))
			w.WriteHeader(http.StatusOK)
			w.Write(data)
		}
	}
	pattern := path[len("html"):]
	return pattern, handler
}

func Start() error {
	r := mux.NewRouter()

	err := fs.WalkDir(htmlFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		r.HandleFunc(htmlHandler(htmlFS, path))
		return nil
	})

	if err != nil {
		return err
	}

	http.Handle("/", r)

	done := make(chan bool)
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", DEFAULT_PORT), nil)
		if err != nil {
			fmt.Println("[ERROR] " + err.Error())
		}
		done <- true
	}()

	<-done

	return nil
}
