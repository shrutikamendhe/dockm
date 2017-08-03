package handler

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/shrutikamendhe/dockm/api"
	httperror "github.com/shrutikamendhe/dockm/api/http/error"
	"github.com/shrutikamendhe/dockm/api/http/security"
)

// TemplatesHandler represents an HTTP API handler for managing templates.
type TemplatesHandler struct {
	*mux.Router
	Logger          *log.Logger
	SettingsService dockm.SettingsService
}

const (
	containerTemplatesURLLinuxServerIo = "http://tools.linuxserver.io/dockm.json"
)

// NewTemplatesHandler returns a new instance of TemplatesHandler.
func NewTemplatesHandler(bouncer *security.RequestBouncer) *TemplatesHandler {
	h := &TemplatesHandler{
		Router: mux.NewRouter(),
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}
	h.Handle("/templates",
		bouncer.AuthenticatedAccess(http.HandlerFunc(h.handleGetTemplates)))
	return h
}

// handleGetTemplates handles GET requests on /templates?key=<key>
func (handler *TemplatesHandler) handleGetTemplates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httperror.WriteMethodNotAllowedResponse(w, []string{http.MethodGet})
		return
	}

	key := r.FormValue("key")
	if key == "" {
		httperror.WriteErrorResponse(w, ErrInvalidQueryFormat, http.StatusBadRequest, handler.Logger)
		return
	}

	var templatesURL string
	if key == "containers" {
		settings, err := handler.SettingsService.Settings()
		if err != nil {
			httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
			return
		}
		templatesURL = settings.TemplatesURL
	} else if key == "linuxserver.io" {
		templatesURL = containerTemplatesURLLinuxServerIo
	} else {
		httperror.WriteErrorResponse(w, ErrInvalidQueryFormat, http.StatusBadRequest, handler.Logger)
		return
	}

	resp, err := http.Get(templatesURL)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
