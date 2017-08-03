package handler

import (
	"github.com/shrutikamendhe/dockm/api"
	httperror "github.com/shrutikamendhe/dockm/api/http/error"
	"github.com/shrutikamendhe/dockm/api/http/security"

	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// UploadHandler represents an HTTP API handler for managing file uploads.
type UploadHandler struct {
	*mux.Router
	Logger      *log.Logger
	FileService dockm.FileService
}

// NewUploadHandler returns a new instance of UploadHandler.
func NewUploadHandler(bouncer *security.RequestBouncer) *UploadHandler {
	h := &UploadHandler{
		Router: mux.NewRouter(),
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}
	h.Handle("/upload/tls/{endpointID}/{certificate:(?:ca|cert|key)}",
		bouncer.AuthenticatedAccess(http.HandlerFunc(h.handlePostUploadTLS)))
	return h
}

func (handler *UploadHandler) handlePostUploadTLS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httperror.WriteMethodNotAllowedResponse(w, []string{http.MethodPost})
		return
	}

	vars := mux.Vars(r)
	endpointID := vars["endpointID"]
	certificate := vars["certificate"]
	ID, err := strconv.Atoi(endpointID)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	file, _, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	var fileType dockm.TLSFileType
	switch certificate {
	case "ca":
		fileType = dockm.TLSFileCA
	case "cert":
		fileType = dockm.TLSFileCert
	case "key":
		fileType = dockm.TLSFileKey
	default:
		httperror.WriteErrorResponse(w, dockm.ErrUndefinedTLSFileType, http.StatusInternalServerError, handler.Logger)
		return
	}

	err = handler.FileService.StoreTLSFile(dockm.EndpointID(ID), fileType, file)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}
}
