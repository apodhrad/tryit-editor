package server

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/apodhrad/tryit-editor/service"
	"github.com/gorilla/mux"
)

const DEFAULT_PORT string = "8080"

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

func htmlHandleFunc(path string) func(w http.ResponseWriter, r *http.Request) {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		data, ok := htmlContent[path]
		if !ok {
			msg := fmt.Sprintf("File '%v' not found", path)
			fmt.Println("[ERROR] " + msg)
			w.Header().Set("Content-Type", CONTENT_TYPE_TEXT)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(msg))
		} else {
			w.Header().Set("Content-Type", contentType(path))
			w.WriteHeader(http.StatusOK)
			w.Write(data)
		}
	}
	return handleFunc
}

func serviceHandler(svc service.Service) (string, func(w http.ResponseWriter, r *http.Request)) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		data, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("[ERROR] " + err.Error())
			w.Header().Set("Content-Type", CONTENT_TYPE_TEXT)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		inputFile, err := os.CreateTemp("", svc.Name()+"-*-input")
		if err != nil {
			fmt.Println("[ERROR] " + err.Error())
			w.Header().Set("Content-Type", CONTENT_TYPE_TEXT)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		_, err = inputFile.Write(data)
		if err != nil {
			fmt.Println("[ERROR] " + err.Error())
			w.Header().Set("Content-Type", CONTENT_TYPE_TEXT)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		out, err := svc.Run(inputFile.Name())
		if err != nil {
			fmt.Println("[ERROR] " + err.Error())
			w.Header().Set("Content-Type", CONTENT_TYPE_TEXT)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		outputFile := strings.ReplaceAll(inputFile.Name(), "input", "output")
		err = os.WriteFile(outputFile, out, 0644)
		if err != nil {
			fmt.Println("[ERROR] " + err.Error())
			w.Header().Set("Content-Type", CONTENT_TYPE_TEXT)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Set("Content-Type", CONTENT_TYPE_TEXT)
		w.WriteHeader(http.StatusOK)
		w.Write(out)
	}
	return "/service/" + svc.Name(), handler
}

var htmlContent map[string][]byte

func registerHtml(r *mux.Router, fsys fs.FS) error {
	htmlContent = make(map[string][]byte)

	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}
		path = path[len("html"):]
		htmlContent[path] = data
		handleFunc := htmlHandleFunc(path)
		r.HandleFunc(path, handleFunc)
		if path == "/index.html" {
			r.HandleFunc("/", handleFunc)
		}
		return nil
	})

	return err
}

func registerService(r *mux.Router, svc service.Service) error {
	pattern, handleFunc := serviceHandler(svc)
	r.HandleFunc(pattern, handleFunc)
	return nil
}

func registerServices(r *mux.Router, svcs []service.Service) error {
	if len(svcs) == 0 {
		return errors.New("At least one service is required!")
	}
	svcOptions := ""
	for _, svc := range svcs {
		err := registerService(r, svc)
		if err != nil {
			return err
		}
		svcOptions += fmt.Sprintf("<option>%v</option>", svc.Name())
	}
	indexData := htmlContent["/index.html"]
	indexData = bytes.ReplaceAll(indexData, []byte("${SERVICE_OPTIONS}"), []byte(svcOptions))
	htmlContent["/index.html"] = indexData
	return nil
}

var server *http.Server
var ctx context.Context
var cancel context.CancelFunc

func Start(svcs []service.Service) (context.Context, error) {
	r := mux.NewRouter()

	err := registerHtml(r, htmlFS)
	if err != nil {
		return ctx, err
	}

	err = registerServices(r, svcs)
	if err != nil {
		return ctx, err
	}

	ctx, cancel = context.WithCancel(context.Background())
	server = &http.Server{
		Addr:    ":" + DEFAULT_PORT,
		Handler: r,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			fmt.Println("[ERROR] " + err.Error())
		}
		cancel()
	}()

	return ctx, nil
}

func Stop() error {
	err := server.Shutdown(ctx)
	if err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}
