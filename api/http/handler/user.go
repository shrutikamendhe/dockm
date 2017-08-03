package handler

import (
	"strconv"
	"strings"

	"github.com/shrutikamendhe/dockm/api"
	httperror "github.com/shrutikamendhe/dockm/api/http/error"
	"github.com/shrutikamendhe/dockm/api/http/security"

	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

// UserHandler represents an HTTP API handler for managing users.
type UserHandler struct {
	*mux.Router
	Logger                 *log.Logger
	UserService            dockm.UserService
	TeamService            dockm.TeamService
	TeamMembershipService  dockm.TeamMembershipService
	ResourceControlService dockm.ResourceControlService
	CryptoService          dockm.CryptoService
}

// NewUserHandler returns a new instance of UserHandler.
func NewUserHandler(bouncer *security.RequestBouncer) *UserHandler {
	h := &UserHandler{
		Router: mux.NewRouter(),
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}
	h.Handle("/users",
		bouncer.RestrictedAccess(http.HandlerFunc(h.handlePostUsers))).Methods(http.MethodPost)
	h.Handle("/users",
		bouncer.RestrictedAccess(http.HandlerFunc(h.handleGetUsers))).Methods(http.MethodGet)
	h.Handle("/users/{id}",
		bouncer.AdministratorAccess(http.HandlerFunc(h.handleGetUser))).Methods(http.MethodGet)
	h.Handle("/users/{id}",
		bouncer.AuthenticatedAccess(http.HandlerFunc(h.handlePutUser))).Methods(http.MethodPut)
	h.Handle("/users/{id}",
		bouncer.AdministratorAccess(http.HandlerFunc(h.handleDeleteUser))).Methods(http.MethodDelete)
	h.Handle("/users/{id}/memberships",
		bouncer.AuthenticatedAccess(http.HandlerFunc(h.handleGetMemberships))).Methods(http.MethodGet)
	h.Handle("/users/{id}/teams",
		bouncer.RestrictedAccess(http.HandlerFunc(h.handleGetTeams))).Methods(http.MethodGet)
	h.Handle("/users/{id}/passwd",
		bouncer.AuthenticatedAccess(http.HandlerFunc(h.handlePostUserPasswd)))
	h.Handle("/users/admin/check",
		bouncer.PublicAccess(http.HandlerFunc(h.handleGetAdminCheck)))
	h.Handle("/users/admin/init",
		bouncer.PublicAccess(http.HandlerFunc(h.handlePostAdminInit)))

	return h
}

// handlePostUsers handles POST requests on /users
func (handler *UserHandler) handlePostUsers(w http.ResponseWriter, r *http.Request) {
	var req postUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.WriteErrorResponse(w, ErrInvalidJSON, http.StatusBadRequest, handler.Logger)
		return
	}

	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		httperror.WriteErrorResponse(w, ErrInvalidRequestFormat, http.StatusBadRequest, handler.Logger)
		return
	}

	securityContext, err := security.RetrieveRestrictedRequestContext(r)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	if !securityContext.IsAdmin && !securityContext.IsTeamLeader {
		httperror.WriteErrorResponse(w, dockm.ErrResourceAccessDenied, http.StatusForbidden, nil)
		return
	}

	if securityContext.IsTeamLeader && req.Role == 1 {
		httperror.WriteErrorResponse(w, dockm.ErrResourceAccessDenied, http.StatusForbidden, nil)
		return
	}

	if strings.ContainsAny(req.Username, " ") {
		httperror.WriteErrorResponse(w, dockm.ErrInvalidUsername, http.StatusBadRequest, handler.Logger)
		return
	}

	var role dockm.UserRole
	if req.Role == 1 {
		role = dockm.AdministratorRole
	} else {
		role = dockm.StandardUserRole
	}

	user, err := handler.UserService.UserByUsername(req.Username)
	if err != nil && err != dockm.ErrUserNotFound {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}
	if user != nil {
		httperror.WriteErrorResponse(w, dockm.ErrUserAlreadyExists, http.StatusConflict, handler.Logger)
		return
	}

	user = &dockm.User{
		Username: req.Username,
		Role:     role,
	}
	user.Password, err = handler.CryptoService.Hash(req.Password)
	if err != nil {
		httperror.WriteErrorResponse(w, dockm.ErrCryptoHashFailure, http.StatusBadRequest, handler.Logger)
		return
	}

	err = handler.UserService.CreateUser(user)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	encodeJSON(w, &postUsersResponse{ID: int(user.ID)}, handler.Logger)
}

type postUsersResponse struct {
	ID int `json:"Id"`
}

type postUsersRequest struct {
	Username string `valid:"required"`
	Password string `valid:"required"`
	Role     int    `valid:"required"`
}

// handleGetUsers handles GET requests on /users
func (handler *UserHandler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	securityContext, err := security.RetrieveRestrictedRequestContext(r)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	users, err := handler.UserService.Users()
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	filteredUsers := security.FilterUsers(users, securityContext)

	for i := range filteredUsers {
		filteredUsers[i].Password = ""
	}

	encodeJSON(w, filteredUsers, handler.Logger)
}

// handlePostUserPasswd handles POST requests on /users/:id/passwd
func (handler *UserHandler) handlePostUserPasswd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httperror.WriteMethodNotAllowedResponse(w, []string{http.MethodPost})
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	userID, err := strconv.Atoi(id)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusBadRequest, handler.Logger)
		return
	}

	var req postUserPasswdRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.WriteErrorResponse(w, ErrInvalidJSON, http.StatusBadRequest, handler.Logger)
		return
	}

	_, err = govalidator.ValidateStruct(req)
	if err != nil {
		httperror.WriteErrorResponse(w, ErrInvalidRequestFormat, http.StatusBadRequest, handler.Logger)
		return
	}

	var password = req.Password

	u, err := handler.UserService.User(dockm.UserID(userID))
	if err == dockm.ErrUserNotFound {
		httperror.WriteErrorResponse(w, err, http.StatusNotFound, handler.Logger)
		return
	} else if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	valid := true
	err = handler.CryptoService.CompareHashAndData(u.Password, password)
	if err != nil {
		valid = false
	}

	encodeJSON(w, &postUserPasswdResponse{Valid: valid}, handler.Logger)
}

type postUserPasswdRequest struct {
	Password string `valid:"required"`
}

type postUserPasswdResponse struct {
	Valid bool `json:"valid"`
}

// handleGetUser handles GET requests on /users/:id
func (handler *UserHandler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	userID, err := strconv.Atoi(id)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusBadRequest, handler.Logger)
		return
	}

	user, err := handler.UserService.User(dockm.UserID(userID))
	if err == dockm.ErrUserNotFound {
		httperror.WriteErrorResponse(w, err, http.StatusNotFound, handler.Logger)
		return
	} else if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	user.Password = ""
	encodeJSON(w, &user, handler.Logger)
}

// handlePutUser handles PUT requests on /users/:id
func (handler *UserHandler) handlePutUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	userID, err := strconv.Atoi(id)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusBadRequest, handler.Logger)
		return
	}

	tokenData, err := security.RetrieveTokenData(r)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	if tokenData.Role != dockm.AdministratorRole && tokenData.ID != dockm.UserID(userID) {
		httperror.WriteErrorResponse(w, dockm.ErrUnauthorized, http.StatusForbidden, handler.Logger)
		return
	}

	var req putUserRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.WriteErrorResponse(w, ErrInvalidJSON, http.StatusBadRequest, handler.Logger)
		return
	}

	_, err = govalidator.ValidateStruct(req)
	if err != nil {
		httperror.WriteErrorResponse(w, ErrInvalidRequestFormat, http.StatusBadRequest, handler.Logger)
		return
	}

	if req.Password == "" && req.Role == 0 {
		httperror.WriteErrorResponse(w, ErrInvalidRequestFormat, http.StatusBadRequest, handler.Logger)
		return
	}

	user, err := handler.UserService.User(dockm.UserID(userID))
	if err == dockm.ErrUserNotFound {
		httperror.WriteErrorResponse(w, err, http.StatusNotFound, handler.Logger)
		return
	} else if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	if req.Password != "" {
		user.Password, err = handler.CryptoService.Hash(req.Password)
		if err != nil {
			httperror.WriteErrorResponse(w, dockm.ErrCryptoHashFailure, http.StatusBadRequest, handler.Logger)
			return
		}
	}

	if req.Role != 0 {
		if tokenData.Role != dockm.AdministratorRole {
			httperror.WriteErrorResponse(w, dockm.ErrUnauthorized, http.StatusForbidden, handler.Logger)
			return
		}
		if req.Role == 1 {
			user.Role = dockm.AdministratorRole
		} else {
			user.Role = dockm.StandardUserRole
		}
	}

	err = handler.UserService.UpdateUser(user.ID, user)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}
}

type putUserRequest struct {
	Password string `valid:"-"`
	Role     int    `valid:"-"`
}

// handlePostAdminInit handles GET requests on /users/admin/check
func (handler *UserHandler) handleGetAdminCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httperror.WriteMethodNotAllowedResponse(w, []string{http.MethodGet})
		return
	}

	users, err := handler.UserService.UsersByRole(dockm.AdministratorRole)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}
	if len(users) == 0 {
		httperror.WriteErrorResponse(w, dockm.ErrUserNotFound, http.StatusNotFound, handler.Logger)
		return
	}
}

// handlePostAdminInit handles POST requests on /users/admin/init
func (handler *UserHandler) handlePostAdminInit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httperror.WriteMethodNotAllowedResponse(w, []string{http.MethodPost})
		return
	}

	var req postAdminInitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.WriteErrorResponse(w, ErrInvalidJSON, http.StatusBadRequest, handler.Logger)
		return
	}

	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		httperror.WriteErrorResponse(w, ErrInvalidRequestFormat, http.StatusBadRequest, handler.Logger)
		return
	}

	user, err := handler.UserService.UserByUsername("admin")
	if err == dockm.ErrUserNotFound {
		user := &dockm.User{
			Username: "admin",
			Role:     dockm.AdministratorRole,
		}
		user.Password, err = handler.CryptoService.Hash(req.Password)
		if err != nil {
			httperror.WriteErrorResponse(w, dockm.ErrCryptoHashFailure, http.StatusBadRequest, handler.Logger)
			return
		}

		err = handler.UserService.CreateUser(user)
		if err != nil {
			httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
			return
		}
	} else if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}
	if user != nil {
		httperror.WriteErrorResponse(w, dockm.ErrAdminAlreadyInitialized, http.StatusForbidden, handler.Logger)
		return
	}
}

type postAdminInitRequest struct {
	Password string `valid:"required"`
}

// handleDeleteUser handles DELETE requests on /users/:id
func (handler *UserHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	userID, err := strconv.Atoi(id)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusBadRequest, handler.Logger)
		return
	}

	_, err = handler.UserService.User(dockm.UserID(userID))

	if err == dockm.ErrUserNotFound {
		httperror.WriteErrorResponse(w, err, http.StatusNotFound, handler.Logger)
		return
	} else if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	err = handler.UserService.DeleteUser(dockm.UserID(userID))
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	err = handler.TeamMembershipService.DeleteTeamMembershipByUserID(dockm.UserID(userID))
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}
}

// handleGetMemberships handles GET requests on /users/:id/memberships
func (handler *UserHandler) handleGetMemberships(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	userID, err := strconv.Atoi(id)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusBadRequest, handler.Logger)
		return
	}

	tokenData, err := security.RetrieveTokenData(r)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	if tokenData.Role != dockm.AdministratorRole && tokenData.ID != dockm.UserID(userID) {
		httperror.WriteErrorResponse(w, dockm.ErrUnauthorized, http.StatusForbidden, handler.Logger)
		return
	}

	memberships, err := handler.TeamMembershipService.TeamMembershipsByUserID(dockm.UserID(userID))
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	encodeJSON(w, memberships, handler.Logger)
}

// handleGetTeams handles GET requests on /users/:id/teams
func (handler *UserHandler) handleGetTeams(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	uid, err := strconv.Atoi(id)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusBadRequest, handler.Logger)
		return
	}
	userID := dockm.UserID(uid)

	securityContext, err := security.RetrieveRestrictedRequestContext(r)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	if !security.AuthorizedUserManagement(userID, securityContext) {
		httperror.WriteErrorResponse(w, dockm.ErrResourceAccessDenied, http.StatusForbidden, handler.Logger)
		return
	}

	teams, err := handler.TeamService.Teams()
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	filteredTeams := security.FilterUserTeams(teams, securityContext)

	encodeJSON(w, filteredTeams, handler.Logger)
}
