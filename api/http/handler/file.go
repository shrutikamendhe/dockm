package handler

import (
	"os"

	"github.com/shrutikamendhe/dockm/api"
	httperror "github.com/shrutikamendhe/dockm/api/http/error"

	"log"
	"net/http"
	"path"
	"strings"
)

// FileHandler represents an HTTP API handler for managing static files.
type FileHandler struct {
	http.Handler
	Logger             *log.Logger
	allowedDirectories map[string]bool
}

// NewFileHandler returns a new instance of FileHandler.
func NewFileHandler(assetPath string) *FileHandler {
	h := &FileHandler{
		Handler: http.FileServer(http.Dir(assetPath)),
		Logger:  log.New(os.Stderr, "", log.LstdFlags),
		allowedDirectories: map[string]bool{
			"/":       true,
			"/css":    true,
			"/js":     true,
			"/images": true,
			"/fonts":  true,
		},
	}
	return h
}

func isHTML(acceptContent []string) bool {
	for _, accept := range acceptContent {
		if strings.Contains(accept, "text/html") {
			return true
		}
	}
	return false
}

func (handler *FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestDirectory := path.Dir(r.URL.Path)
	if !handler.allowedDirectories[requestDirectory] {
		httperror.WriteErrorResponse(w, dockm.ErrResourceNotFound, http.StatusNotFound, handler.Logger)
		return
	}

	if !isHTML(r.Header["Accept"]) {
		w.Header().Set("Cache-Control", "max-age=31536000")
	} else {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	}

	handler.Handler.ServeHTTP(w, r)
}
