package proxy

import (
	"net/http"

	"github.com/shrutikamendhe/dockm/api"
)

const (
	// ErrDockerContainerIdentifierNotFound defines an error raised when DockM is unable to find a container identifier
	ErrDockerContainerIdentifierNotFound = dockm.Error("Docker container identifier not found")
	containerIdentifier                  = "Id"
	containerLabelForServiceIdentifier   = "com.docker.swarm.service.id"
)

// containerListOperation extracts the response as a JSON object, loop through the containers array
// decorate and/or filter the containers based on resource controls before rewriting the response
func containerListOperation(request *http.Request, response *http.Response, executor *operationExecutor) error {
	var err error
	// ContainerList response is a JSON array
	// https://docs.docker.com/engine/api/v1.28/#operation/ContainerList
	responseArray, err := getResponseAsJSONArray(response)
	if err != nil {
		return err
	}

	if executor.operationContext.isAdmin {
		responseArray, err = decorateContainerList(responseArray, executor.operationContext.resourceControls)
	} else {
		responseArray, err = filterContainerList(responseArray, executor.operationContext.resourceControls,
			executor.operationContext.userID, executor.operationContext.userTeamIDs)
	}
	if err != nil {
		return err
	}

	if executor.labelBlackList != nil {
		responseArray, err = filterContainersWithBlackListedLabels(responseArray, executor.labelBlackList)
		if err != nil {
			return err
		}
	}

	return rewriteResponse(response, responseArray, http.StatusOK)
}

// containerInspectOperation extracts the response as a JSON object, verify that the user
// has access to the container based on resource control (check are done based on the containerID and optional Swarm service ID)
// and either rewrite an access denied response or a decorated container.
func containerInspectOperation(request *http.Request, response *http.Response, executor *operationExecutor) error {
	// ContainerInspect response is a JSON object
	// https://docs.docker.com/engine/api/v1.28/#operation/ContainerInspect
	responseObject, err := getResponseAsJSONOBject(response)
	if err != nil {
		return err
	}

	if responseObject[containerIdentifier] == nil {
		return ErrDockerContainerIdentifierNotFound
	}
	containerID := responseObject[containerIdentifier].(string)

	resourceControl := getResourceControlByResourceID(containerID, executor.operationContext.resourceControls)
	if resourceControl != nil {
		if executor.operationContext.isAdmin || canUserAccessResource(executor.operationContext.userID,
			executor.operationContext.userTeamIDs, resourceControl) {
			responseObject = decorateObject(responseObject, resourceControl)
		} else {
			return rewriteAccessDeniedResponse(response)
		}
	}

	containerLabels := extractContainerLabelsFromContainerInspectObject(responseObject)
	if containerLabels != nil && containerLabels[containerLabelForServiceIdentifier] != nil {
		serviceID := containerLabels[containerLabelForServiceIdentifier].(string)
		resourceControl := getResourceControlByResourceID(serviceID, executor.operationContext.resourceControls)
		if resourceControl != nil {
			if executor.operationContext.isAdmin || canUserAccessResource(executor.operationContext.userID,
				executor.operationContext.userTeamIDs, resourceControl) {
				responseObject = decorateObject(responseObject, resourceControl)
			} else {
				return rewriteAccessDeniedResponse(response)
			}
		}
	}

	return rewriteResponse(response, responseObject, http.StatusOK)
}

// extractContainerLabelsFromContainerInspectObject retrieve the Labels of the container if present.
// Container schema reference: https://docs.docker.com/engine/api/v1.28/#operation/ContainerInspect
func extractContainerLabelsFromContainerInspectObject(responseObject map[string]interface{}) map[string]interface{} {
	// Labels are stored under Config.Labels
	containerConfigObject := extractJSONField(responseObject, "Config")
	if containerConfigObject != nil {
		containerLabelsObject := extractJSONField(containerConfigObject, "Labels")
		return containerLabelsObject
	}
	return nil
}

// extractContainerLabelsFromContainerListObject retrieve the Labels of the container if present.
// Container schema reference: https://docs.docker.com/engine/api/v1.28/#operation/ContainerList
func extractContainerLabelsFromContainerListObject(responseObject map[string]interface{}) map[string]interface{} {
	// Labels are stored under Labels
	containerLabelsObject := extractJSONField(responseObject, "Labels")
	return containerLabelsObject
}
