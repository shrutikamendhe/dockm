package handler

import (
	"strconv"

	"github.com/shrutikamendhe/dockm/api"
	httperror "github.com/shrutikamendhe/dockm/api/http/error"
	"github.com/shrutikamendhe/dockm/api/http/proxy"
	"github.com/shrutikamendhe/dockm/api/http/security"

	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// DockerHandler represents an HTTP API handler for proxying requests to the Docker API.
type DockerHandler struct {
	*mux.Router
	Logger                *log.Logger
	EndpointService       dockm.EndpointService
	TeamMembershipService dockm.TeamMembershipService
	ProxyManager          *proxy.Manager
}

// NewDockerHandler returns a new instance of DockerHandler.
func NewDockerHandler(bouncer *security.RequestBouncer) *DockerHandler {
	h := &DockerHandler{
		Router: mux.NewRouter(),
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}
	h.PathPrefix("/{id}/docker").Handler(
		bouncer.AuthenticatedAccess(http.HandlerFunc(h.proxyRequestsToDockerAPI)))
	return h
}

func (handler *DockerHandler) checkEndpointAccessControl(endpoint *dockm.Endpoint, userID dockm.UserID) bool {
	for _, authorizedUserID := range endpoint.AuthorizedUsers {
		if authorizedUserID == userID {
			return true
		}
	}

	memberships, _ := handler.TeamMembershipService.TeamMembershipsByUserID(userID)
	for _, authorizedTeamID := range endpoint.AuthorizedTeams {
		for _, membership := range memberships {
			if membership.TeamID == authorizedTeamID {
				return true
			}
		}
	}
	return false
}

func (handler *DockerHandler) proxyRequestsToDockerAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	parsedID, err := strconv.Atoi(id)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusBadRequest, handler.Logger)
		return
	}

	endpointID := dockm.EndpointID(parsedID)
	endpoint, err := handler.EndpointService.Endpoint(endpointID)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	tokenData, err := security.RetrieveTokenData(r)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}
	if tokenData.Role != dockm.AdministratorRole && !handler.checkEndpointAccessControl(endpoint, tokenData.ID) {
		httperror.WriteErrorResponse(w, dockm.ErrEndpointAccessDenied, http.StatusForbidden, handler.Logger)
		return
	}

	var proxy http.Handler
	proxy = handler.ProxyManager.GetProxy(string(endpointID))
	if proxy == nil {
		proxy, err = handler.ProxyManager.CreateAndRegisterProxy(endpoint)
		if err != nil {
			httperror.WriteErrorResponse(w, err, http.StatusBadRequest, handler.Logger)
			return
		}
	}

	http.StripPrefix("/"+id+"/docker", proxy).ServeHTTP(w, r)
}
