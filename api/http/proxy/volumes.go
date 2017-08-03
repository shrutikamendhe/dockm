package proxy

import (
	"net/http"

	"github.com/shrutikamendhe/dockm/api"
)

const (
	// ErrDockerVolumeIdentifierNotFound defines an error raised when DockM is unable to find a volume identifier
	ErrDockerVolumeIdentifierNotFound = dockm.Error("Docker volume identifier not found")
	volumeIdentifier                  = "Name"
)

// volumeListOperation extracts the response as a JSON object, loop through the volume array
// decorate and/or filter the volumes based on resource controls before rewriting the response
func volumeListOperation(request *http.Request, response *http.Response, executor *operationExecutor) error {
	var err error
	// VolumeList response is a JSON object
	// https://docs.docker.com/engine/api/v1.28/#operation/VolumeList
	responseObject, err := getResponseAsJSONOBject(response)
	if err != nil {
		return err
	}

	// The "Volumes" field contains the list of volumes as an array of JSON objects
	// Response schema reference: https://docs.docker.com/engine/api/v1.28/#operation/VolumeList
	if responseObject["Volumes"] != nil {
		volumeData := responseObject["Volumes"].([]interface{})

		if executor.operationContext.isAdmin {
			volumeData, err = decorateVolumeList(volumeData, executor.operationContext.resourceControls)
		} else {
			volumeData, err = filterVolumeList(volumeData, executor.operationContext.resourceControls, executor.operationContext.userID, executor.operationContext.userTeamIDs)
		}
		if err != nil {
			return err
		}

		// Overwrite the original volume list
		responseObject["Volumes"] = volumeData
	}

	return rewriteResponse(response, responseObject, http.StatusOK)
}

// volumeInspectOperation extracts the response as a JSON object, verify that the user
// has access to the volume based on resource control and either rewrite an access denied response
// or a decorated volume.
func volumeInspectOperation(request *http.Request, response *http.Response, executor *operationExecutor) error {
	// VolumeInspect response is a JSON object
	// https://docs.docker.com/engine/api/v1.28/#operation/VolumeInspect
	responseObject, err := getResponseAsJSONOBject(response)
	if err != nil {
		return err
	}

	if responseObject[volumeIdentifier] == nil {
		return ErrDockerVolumeIdentifierNotFound
	}
	volumeID := responseObject[volumeIdentifier].(string)

	resourceControl := getResourceControlByResourceID(volumeID, executor.operationContext.resourceControls)
	if resourceControl != nil {
		if executor.operationContext.isAdmin || canUserAccessResource(executor.operationContext.userID, executor.operationContext.userTeamIDs, resourceControl) {
			responseObject = decorateObject(responseObject, resourceControl)
		} else {
			return rewriteAccessDeniedResponse(response)
		}
	}

	return rewriteResponse(response, responseObject, http.StatusOK)
}
