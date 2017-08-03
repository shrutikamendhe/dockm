package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/shrutikamendhe/dockm/api"
	httperror "github.com/shrutikamendhe/dockm/api/http/error"
)

// Handler is a collection of all the service handlers.
type Handler struct {
	AuthHandler           *AuthHandler
	UserHandler           *UserHandler
	TeamHandler           *TeamHandler
	TeamMembershipHandler *TeamMembershipHandler
	EndpointHandler       *EndpointHandler
	RegistryHandler       *RegistryHandler
	DockerHubHandler      *DockerHubHandler
	ResourceHandler       *ResourceHandler
	StackHandler          *StackHandler
	StatusHandler         *StatusHandler
	SettingsHandler       *SettingsHandler
	TemplatesHandler      *TemplatesHandler
	DockerHandler         *DockerHandler
	WebSocketHandler      *WebSocketHandler
	UploadHandler         *UploadHandler
	FileHandler           *FileHandler
}

const (
	// ErrInvalidJSON defines an error raised the app is unable to parse request data
	ErrInvalidJSON = dockm.Error("Invalid JSON")
	// ErrInvalidRequestFormat defines an error raised when the format of the data sent in a request is not valid
	ErrInvalidRequestFormat = dockm.Error("Invalid request data format")
	// ErrInvalidQueryFormat defines an error raised when the data sent in the query or the URL is invalid
	ErrInvalidQueryFormat = dockm.Error("Invalid query format")
	// ErrEmptyResponseBody defines an error raised when dockm excepts to parse the body of a HTTP response and there is nothing to parse
	// ErrEmptyResponseBody = dockm.Error("Empty response body")
)

// ServeHTTP delegates a request to the appropriate subhandler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/auth") {
		http.StripPrefix("/api", h.AuthHandler).ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/api/users") {
		http.StripPrefix("/api", h.UserHandler).ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/api/teams") {
		http.StripPrefix("/api", h.TeamHandler).ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/api/team_memberships") {
		http.StripPrefix("/api", h.TeamMembershipHandler).ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/api/endpoints") {

		if strings.Contains(r.URL.Path, "stacks") {
			http.StripPrefix("/api/endpoints", h.StackHandler).ServeHTTP(w, r)
    } else if strings.Contains(r.URL.Path, "/docker") {
			http.StripPrefix("/api/endpoints", h.DockerHandler).ServeHTTP(w, r)
		} else {
			http.StripPrefix("/api", h.EndpointHandler).ServeHTTP(w, r)
		}
	} else if strings.HasPrefix(r.URL.Path, "/api/registries") {
		http.StripPrefix("/api", h.RegistryHandler).ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/api/dockerhub") {
		http.StripPrefix("/api", h.DockerHubHandler).ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/api/resource_controls") {
		http.StripPrefix("/api", h.ResourceHandler).ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/api/settings") {
		http.StripPrefix("/api", h.SettingsHandler).ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/api/status") {
		http.StripPrefix("/api", h.StatusHandler).ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/api/templates") {
		http.StripPrefix("/api", h.TemplatesHandler).ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/api/upload") {
		http.StripPrefix("/api", h.UploadHandler).ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/api/websocket") {
		http.StripPrefix("/api", h.WebSocketHandler).ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/") {
		h.FileHandler.ServeHTTP(w, r)
	}
}

// encodeJSON encodes v to w in JSON format. Error() is called if encoding fails.
func encodeJSON(w http.ResponseWriter, v interface{}, logger *log.Logger) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, logger)
	}
}
