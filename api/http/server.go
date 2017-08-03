package http

import (
	"github.com/shrutikamendhe/dockm/api"
	"github.com/shrutikamendhe/dockm/api/http/handler"
	"github.com/shrutikamendhe/dockm/api/http/proxy"
	"github.com/shrutikamendhe/dockm/api/http/security"

	"net/http"
)

// Server implements the dockm.Server interface
type Server struct {
	BindAddress            string
	AssetsPath             string
	AuthDisabled           bool
	EndpointManagement     bool
	Status                 *dockm.Status
	UserService            dockm.UserService
	TeamService            dockm.TeamService
	TeamMembershipService  dockm.TeamMembershipService
	EndpointService        dockm.EndpointService
	ResourceControlService dockm.ResourceControlService
	SettingsService        dockm.SettingsService
	CryptoService          dockm.CryptoService
	JWTService             dockm.JWTService
	FileService            dockm.FileService
	RegistryService        dockm.RegistryService
	DockerHubService       dockm.DockerHubService
	StackService           dockm.StackService
	StackManager           dockm.StackManager
	Handler                *handler.Handler
	SSL                    bool
	SSLCert                string
	SSLKey                 string
}

// Start starts the HTTP server
func (server *Server) Start() error {
	requestBouncer := security.NewRequestBouncer(server.JWTService, server.TeamMembershipService, server.AuthDisabled)
	proxyManager := proxy.NewManager(server.ResourceControlService, server.TeamMembershipService, server.SettingsService)

	var fileHandler = handler.NewFileHandler(server.AssetsPath)
	var authHandler = handler.NewAuthHandler(requestBouncer, server.AuthDisabled)
	authHandler.UserService = server.UserService
	authHandler.CryptoService = server.CryptoService
	authHandler.JWTService = server.JWTService
	var userHandler = handler.NewUserHandler(requestBouncer)
	userHandler.UserService = server.UserService
	userHandler.TeamService = server.TeamService
	userHandler.TeamMembershipService = server.TeamMembershipService
	userHandler.CryptoService = server.CryptoService
	userHandler.ResourceControlService = server.ResourceControlService
	var teamHandler = handler.NewTeamHandler(requestBouncer)
	teamHandler.TeamService = server.TeamService
	teamHandler.TeamMembershipService = server.TeamMembershipService
	var teamMembershipHandler = handler.NewTeamMembershipHandler(requestBouncer)
	teamMembershipHandler.TeamMembershipService = server.TeamMembershipService
	var statusHandler = handler.NewStatusHandler(requestBouncer, server.Status)
	var settingsHandler = handler.NewSettingsHandler(requestBouncer)
	settingsHandler.SettingsService = server.SettingsService
	var templatesHandler = handler.NewTemplatesHandler(requestBouncer)
	templatesHandler.SettingsService = server.SettingsService
	var dockerHandler = handler.NewDockerHandler(requestBouncer)
	dockerHandler.EndpointService = server.EndpointService
	dockerHandler.TeamMembershipService = server.TeamMembershipService
	dockerHandler.ProxyManager = proxyManager
	var websocketHandler = handler.NewWebSocketHandler()
	websocketHandler.EndpointService = server.EndpointService
	var endpointHandler = handler.NewEndpointHandler(requestBouncer, server.EndpointManagement)
	endpointHandler.EndpointService = server.EndpointService
	endpointHandler.FileService = server.FileService
	endpointHandler.ProxyManager = proxyManager
	var registryHandler = handler.NewRegistryHandler(requestBouncer)
	registryHandler.RegistryService = server.RegistryService
	var dockerHubHandler = handler.NewDockerHubHandler(requestBouncer)
	dockerHubHandler.DockerHubService = server.DockerHubService
	var resourceHandler = handler.NewResourceHandler(requestBouncer)
	resourceHandler.ResourceControlService = server.ResourceControlService
	var uploadHandler = handler.NewUploadHandler(requestBouncer)
	uploadHandler.FileService = server.FileService
	var stackHandler = handler.NewStackHandler(requestBouncer)
	stackHandler.FileService = server.FileService
	stackHandler.StackService = server.StackService
	stackHandler.EndpointService = server.EndpointService
	stackHandler.StackManager = server.StackManager

	server.Handler = &handler.Handler{
		AuthHandler:           authHandler,
		UserHandler:           userHandler,
		TeamHandler:           teamHandler,
		TeamMembershipHandler: teamMembershipHandler,
		EndpointHandler:       endpointHandler,
		RegistryHandler:       registryHandler,
		DockerHubHandler:      dockerHubHandler,
		ResourceHandler:       resourceHandler,
		SettingsHandler:       settingsHandler,
		StatusHandler:         statusHandler,
		StackHandler:          stackHandler,
		TemplatesHandler:      templatesHandler,
		DockerHandler:         dockerHandler,
		WebSocketHandler:      websocketHandler,
		FileHandler:           fileHandler,
		UploadHandler:         uploadHandler,
	}

	if server.SSL {
		return http.ListenAndServeTLS(server.BindAddress, server.SSLCert, server.SSLKey, server.Handler)
	}
	return http.ListenAndServe(server.BindAddress, server.Handler)
}
