package proxy

import "github.com/shrutikamendhe/dockm/api"

// decorateVolumeList loops through all volumes and will decorate any volume with an existing resource control.
// Volume object schema reference: https://docs.docker.com/engine/api/v1.28/#operation/VolumeList
func decorateVolumeList(volumeData []interface{}, resourceControls []dockm.ResourceControl) ([]interface{}, error) {
	decoratedVolumeData := make([]interface{}, 0)

	for _, volume := range volumeData {

		volumeObject := volume.(map[string]interface{})
		if volumeObject[volumeIdentifier] == nil {
			return nil, ErrDockerVolumeIdentifierNotFound
		}

		volumeID := volumeObject[volumeIdentifier].(string)
		resourceControl := getResourceControlByResourceID(volumeID, resourceControls)
		if resourceControl != nil {
			volumeObject = decorateObject(volumeObject, resourceControl)
		}
		decoratedVolumeData = append(decoratedVolumeData, volumeObject)
	}

	return decoratedVolumeData, nil
}

// decorateContainerList loops through all containers and will decorate any container with an existing resource control.
// Check is based on the container ID and optional Swarm service ID.
// Container object schema reference: https://docs.docker.com/engine/api/v1.28/#operation/ContainerList
func decorateContainerList(containerData []interface{}, resourceControls []dockm.ResourceControl) ([]interface{}, error) {
	decoratedContainerData := make([]interface{}, 0)

	for _, container := range containerData {

		containerObject := container.(map[string]interface{})
		if containerObject[containerIdentifier] == nil {
			return nil, ErrDockerContainerIdentifierNotFound
		}

		containerID := containerObject[containerIdentifier].(string)
		resourceControl := getResourceControlByResourceID(containerID, resourceControls)
		if resourceControl != nil {
			containerObject = decorateObject(containerObject, resourceControl)
		}

		containerLabels := extractContainerLabelsFromContainerListObject(containerObject)
		if containerLabels != nil && containerLabels[containerLabelForServiceIdentifier] != nil {
			serviceID := containerLabels[containerLabelForServiceIdentifier].(string)
			resourceControl := getResourceControlByResourceID(serviceID, resourceControls)
			if resourceControl != nil {
				containerObject = decorateObject(containerObject, resourceControl)
			}
		}

		decoratedContainerData = append(decoratedContainerData, containerObject)
	}

	return decoratedContainerData, nil
}

// decorateServiceList loops through all services and will decorate any service with an existing resource control.
// Service object schema reference: https://docs.docker.com/engine/api/v1.28/#operation/ServiceList
func decorateServiceList(serviceData []interface{}, resourceControls []dockm.ResourceControl) ([]interface{}, error) {
	decoratedServiceData := make([]interface{}, 0)

	for _, service := range serviceData {

		serviceObject := service.(map[string]interface{})
		if serviceObject[serviceIdentifier] == nil {
			return nil, ErrDockerServiceIdentifierNotFound
		}

		serviceID := serviceObject[serviceIdentifier].(string)
		resourceControl := getResourceControlByResourceID(serviceID, resourceControls)
		if resourceControl != nil {
			serviceObject = decorateObject(serviceObject, resourceControl)
		}
		decoratedServiceData = append(decoratedServiceData, serviceObject)
	}

	return decoratedServiceData, nil
}

func decorateObject(object map[string]interface{}, resourceControl *dockm.ResourceControl) map[string]interface{} {
	metadata := make(map[string]interface{})
	metadata["ResourceControl"] = resourceControl
	object["DockM"] = metadata
	return object
}
