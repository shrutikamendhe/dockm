angular.module('dockm')
.controller('porImageRegistryController', ['$q', 'RegistryService', 'DockerHubService', 'Notifications',
function ($q, RegistryService, DockerHubService, Notifications) {
  var ctrl = this;

  function initComponent() {
    $q.all({
      registries: RegistryService.registries(),
      dockerhub: DockerHubService.dockerhub()
    })
    .then(function success(data) {
      var dockerhub = data.dockerhub;
      var registries = data.registries;
      ctrl.availableRegistries = [dockerhub].concat(registries);
      ctrl.registry = dockerhub;
    })
    .catch(function error(err) {
      Notifications.error('Failure', err, 'Unable to retrieve registries');
    });
  }

  initComponent();
}]);
