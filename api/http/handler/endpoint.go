package handler

import (
	"github.com/shrutikamendhe/dockm/api"
	httperror "github.com/shrutikamendhe/dockm/api/http/error"
	"github.com/shrutikamendhe/dockm/api/http/proxy"
	"github.com/shrutikamendhe/dockm/api/http/security"

	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

// EndpointHandler represents an HTTP API handler for managing Docker endpoints.
type EndpointHandler struct {
	*mux.Router
	Logger                      *log.Logger
	authorizeEndpointManagement bool
	EndpointService             dockm.EndpointService
	FileService                 dockm.FileService
	ProxyManager                *proxy.Manager
}

const (
	// ErrEndpointManagementDisabled is an error raised when trying to access the endpoints management endpoints
	// when the server has been started with the --external-endpoints flag
	ErrEndpointManagementDisabled = dockm.Error("Endpoint management is disabled")
)

// NewEndpointHandler returns a new instance of EndpointHandler.
func NewEndpointHandler(bouncer *security.RequestBouncer, authorizeEndpointManagement bool) *EndpointHandler {
	h := &EndpointHandler{
		Router: mux.NewRouter(),
		Logger: log.New(os.Stderr, "", log.LstdFlags),
		authorizeEndpointManagement: authorizeEndpointManagement,
	}
	h.Handle("/endpoints",
		bouncer.AdministratorAccess(http.HandlerFunc(h.handlePostEndpoints))).Methods(http.MethodPost)
	h.Handle("/endpoints",
		bouncer.RestrictedAccess(http.HandlerFunc(h.handleGetEndpoints))).Methods(http.MethodGet)
	h.Handle("/endpoints/{id}",
		bouncer.AdministratorAccess(http.HandlerFunc(h.handleGetEndpoint))).Methods(http.MethodGet)
	h.Handle("/endpoints/{id}",
		bouncer.AdministratorAccess(http.HandlerFunc(h.handlePutEndpoint))).Methods(http.MethodPut)
	h.Handle("/endpoints/{id}/access",
		bouncer.AdministratorAccess(http.HandlerFunc(h.handlePutEndpointAccess))).Methods(http.MethodPut)
	h.Handle("/endpoints/{id}",
		bouncer.AdministratorAccess(http.HandlerFunc(h.handleDeleteEndpoint))).Methods(http.MethodDelete)

	return h
}

// handleGetEndpoints handles GET requests on /endpoints
func (handler *EndpointHandler) handleGetEndpoints(w http.ResponseWriter, r *http.Request) {
	securityContext, err := security.RetrieveRestrictedRequestContext(r)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	endpoints, err := handler.EndpointService.Endpoints()
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	filteredEndpoints, err := security.FilterEndpoints(endpoints, securityContext)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	encodeJSON(w, filteredEndpoints, handler.Logger)
}

// handlePostEndpoints handles POST requests on /endpoints
func (handler *EndpointHandler) handlePostEndpoints(w http.ResponseWriter, r *http.Request) {
	if !handler.authorizeEndpointManagement {
		httperror.WriteErrorResponse(w, ErrEndpointManagementDisabled, http.StatusServiceUnavailable, handler.Logger)
		return
	}

	var req postEndpointsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.WriteErrorResponse(w, ErrInvalidJSON, http.StatusBadRequest, handler.Logger)
		return
	}

	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		httperror.WriteErrorResponse(w, ErrInvalidRequestFormat, http.StatusBadRequest, handler.Logger)
		return
	}

	endpoint := &dockm.Endpoint{
		Name:            req.Name,
		URL:             req.URL,
		PublicURL:       req.PublicURL,
		TLS:             req.TLS,
		AuthorizedUsers: []dockm.UserID{},
		AuthorizedTeams: []dockm.TeamID{},
	}

	err = handler.EndpointService.CreateEndpoint(endpoint)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	if req.TLS {
		caCertPath, _ := handler.FileService.GetPathForTLSFile(endpoint.ID, dockm.TLSFileCA)
		endpoint.TLSCACertPath = caCertPath
		certPath, _ := handler.FileService.GetPathForTLSFile(endpoint.ID, dockm.TLSFileCert)
		endpoint.TLSCertPath = certPath
		keyPath, _ := handler.FileService.GetPathForTLSFile(endpoint.ID, dockm.TLSFileKey)
		endpoint.TLSKeyPath = keyPath
		err = handler.EndpointService.UpdateEndpoint(endpoint.ID, endpoint)
		if err != nil {
			httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
			return
		}
	}

	encodeJSON(w, &postEndpointsResponse{ID: int(endpoint.ID)}, handler.Logger)
}

type postEndpointsRequest struct {
	Name      string `valid:"required"`
	URL       string `valid:"required"`
	PublicURL string `valid:"-"`
	TLS       bool
}

type postEndpointsResponse struct {
	ID int `json:"Id"`
}

// handleGetEndpoint handles GET requests on /endpoints/:id
func (handler *EndpointHandler) handleGetEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	endpointID, err := strconv.Atoi(id)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusBadRequest, handler.Logger)
		return
	}

	endpoint, err := handler.EndpointService.Endpoint(dockm.EndpointID(endpointID))
	if err == dockm.ErrEndpointNotFound {
		httperror.WriteErrorResponse(w, err, http.StatusNotFound, handler.Logger)
		return
	} else if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	encodeJSON(w, endpoint, handler.Logger)
}

// handlePutEndpointAccess handles PUT requests on /endpoints/:id/access
func (handler *EndpointHandler) handlePutEndpointAccess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	endpointID, err := strconv.Atoi(id)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusBadRequest, handler.Logger)
		return
	}

	var req putEndpointAccessRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.WriteErrorResponse(w, ErrInvalidJSON, http.StatusBadRequest, handler.Logger)
		return
	}

	_, err = govalidator.ValidateStruct(req)
	if err != nil {
		httperror.WriteErrorResponse(w, ErrInvalidRequestFormat, http.StatusBadRequest, handler.Logger)
		return
	}

	endpoint, err := handler.EndpointService.Endpoint(dockm.EndpointID(endpointID))
	if err == dockm.ErrEndpointNotFound {
		httperror.WriteErrorResponse(w, err, http.StatusNotFound, handler.Logger)
		return
	} else if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	if req.AuthorizedUsers != nil {
		authorizedUserIDs := []dockm.UserID{}
		for _, value := range req.AuthorizedUsers {
			authorizedUserIDs = append(authorizedUserIDs, dockm.UserID(value))
		}
		endpoint.AuthorizedUsers = authorizedUserIDs
	}

	if req.AuthorizedTeams != nil {
		authorizedTeamIDs := []dockm.TeamID{}
		for _, value := range req.AuthorizedTeams {
			authorizedTeamIDs = append(authorizedTeamIDs, dockm.TeamID(value))
		}
		endpoint.AuthorizedTeams = authorizedTeamIDs
	}

	err = handler.EndpointService.UpdateEndpoint(endpoint.ID, endpoint)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}
}

type putEndpointAccessRequest struct {
	AuthorizedUsers []int `valid:"-"`
	AuthorizedTeams []int `valid:"-"`
}

// handlePutEndpoint handles PUT requests on /endpoints/:id
func (handler *EndpointHandler) handlePutEndpoint(w http.ResponseWriter, r *http.Request) {
	if !handler.authorizeEndpointManagement {
		httperror.WriteErrorResponse(w, ErrEndpointManagementDisabled, http.StatusServiceUnavailable, handler.Logger)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	endpointID, err := strconv.Atoi(id)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusBadRequest, handler.Logger)
		return
	}

	var req putEndpointsRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.WriteErrorResponse(w, ErrInvalidJSON, http.StatusBadRequest, handler.Logger)
		return
	}

	_, err = govalidator.ValidateStruct(req)
	if err != nil {
		httperror.WriteErrorResponse(w, ErrInvalidRequestFormat, http.StatusBadRequest, handler.Logger)
		return
	}

	endpoint, err := handler.EndpointService.Endpoint(dockm.EndpointID(endpointID))
	if err == dockm.ErrEndpointNotFound {
		httperror.WriteErrorResponse(w, err, http.StatusNotFound, handler.Logger)
		return
	} else if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	if req.Name != "" {
		endpoint.Name = req.Name
	}

	if req.URL != "" {
		endpoint.URL = req.URL
	}

	if req.PublicURL != "" {
		endpoint.PublicURL = req.PublicURL
	}

	if req.TLS {
		endpoint.TLS = true
		caCertPath, _ := handler.FileService.GetPathForTLSFile(endpoint.ID, dockm.TLSFileCA)
		endpoint.TLSCACertPath = caCertPath
		certPath, _ := handler.FileService.GetPathForTLSFile(endpoint.ID, dockm.TLSFileCert)
		endpoint.TLSCertPath = certPath
		keyPath, _ := handler.FileService.GetPathForTLSFile(endpoint.ID, dockm.TLSFileKey)
		endpoint.TLSKeyPath = keyPath
	} else {
		endpoint.TLS = false
		endpoint.TLSCACertPath = ""
		endpoint.TLSCertPath = ""
		endpoint.TLSKeyPath = ""
		err = handler.FileService.DeleteTLSFiles(endpoint.ID)
		if err != nil {
			httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
			return
		}
	}

	_, err = handler.ProxyManager.CreateAndRegisterProxy(endpoint)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	err = handler.EndpointService.UpdateEndpoint(endpoint.ID, endpoint)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}
}

type putEndpointsRequest struct {
	Name      string `valid:"-"`
	URL       string `valid:"-"`
	PublicURL string `valid:"-"`
	TLS       bool   `valid:"-"`
}

// handleDeleteEndpoint handles DELETE requests on /endpoints/:id
func (handler *EndpointHandler) handleDeleteEndpoint(w http.ResponseWriter, r *http.Request) {
	if !handler.authorizeEndpointManagement {
		httperror.WriteErrorResponse(w, ErrEndpointManagementDisabled, http.StatusServiceUnavailable, handler.Logger)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	endpointID, err := strconv.Atoi(id)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusBadRequest, handler.Logger)
		return
	}

	endpoint, err := handler.EndpointService.Endpoint(dockm.EndpointID(endpointID))

	if err == dockm.ErrEndpointNotFound {
		httperror.WriteErrorResponse(w, err, http.StatusNotFound, handler.Logger)
		return
	} else if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	handler.ProxyManager.DeleteProxy(string(endpointID))

	err = handler.EndpointService.DeleteEndpoint(dockm.EndpointID(endpointID))
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	if endpoint.TLS {
		err = handler.FileService.DeleteTLSFiles(dockm.EndpointID(endpointID))
		if err != nil {
			httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
			return
		}
	}
}
