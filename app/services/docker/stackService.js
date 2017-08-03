angular.module('dockm.services')
.factory('StackService', ['$q', 'Stack', 'ContainerService', 'ServiceService', 'TaskService', 'StackHelper', function StackServiceFactory($q, Stack, ContainerService, ServiceService, TaskService, StackHelper) {
  'use strict';
  var service = {};

//TODO: remove useless methods

  service.createStack = function(name, composeFile, envFile) {
    return Stack.create({Name: name, ComposeFileContent: composeFile, EnvFileContent: envFile}).$promise;
  };

  service.retrieveStacksAndAnonymousStacks = function(includeServices) {
    var deferred = $q.defer();

    $q.all({
      stacks: service.stacks(),
      discoveredStacks: service.discoverStacks(includeServices)
    })
    .then(function success(data) {
      var stacks = data.stacks;
      var discoveredStacks = data.discoveredStacks;
      var anonymousStacks = StackHelper.mergeStacksAndDiscoveredStacks(stacks, discoveredStacks);
      deferred.resolve({ stacks: stacks, anonymousStacks: anonymousStacks });
    })
    .catch(function error(err) {
      deferred.reject({ msg: 'Unable to retrieve stacks', err: err });
    });

    return deferred.promise;
  };

  service.stack = function(id) {
    var deferred = $q.defer();

    Stack.get({id: id}).$promise
    .then(function success(data) {
      var stack = new StackViewModel(data);
      deferred.resolve(stack);
    })
    .catch(function error(err) {
      deferred.reject({ msg: 'Unable to retrieve stack details', err: err });
    });

    return deferred.promise;
  };

  service.stacks = function() {
    var deferred = $q.defer();

    Stack.query().$promise
    .then(function success(data) {
      var stacks = data.map(function (item) {
        return new StackViewModel(item);
      });
      deferred.resolve(stacks);
    })
    .catch(function error(err) {
      deferred.reject({ msg: 'Unable to retrieve stacks', err: err });
    });

    return deferred.promise;
  };

  service.discoverStacks = function(includeServices) {
    var deferred = $q.defer();

    $q.all({
      containers: ContainerService.containers(1),
      services: includeServices ? ServiceService.services() : []
    })
    .then(function success(data) {
      var containers = data.containers;
      var composeV2Stacks = StackHelper.getComposeV2StacksFromContainers(containers);
      var services = data.services;
      var composeV3Stacks = StackHelper.getComposeV3StacksFromServices(services);

      var stacks = composeV2Stacks.concat(composeV3Stacks);
      deferred.resolve(stacks);
    })
    .catch(function error(err) {
      deferred.reject({ msg: 'Stack discovery failure', err: err });
    });

    return deferred.promise;
  };


  service.getStackV2ServicesAndContainers = function(name) {
    var deferred = $q.defer();

    var filters = {
      label: ['com.docker.compose.project=' + name]
    };

    ContainerService.containers(1, filters)
    .then(function success(data) {
      var containers = data;
      var services = StackHelper.getComposeV2ServicesFromContainers(containers);
      deferred.resolve({services: services, containers: containers});
    })
    .catch(function error(err) {
      deferred.reject({ msg: 'Unable to retrieve stack details', err: err });
    });

    return deferred.promise;
  };

  // service.stackV2 = function(name) {
  //   var deferred = $q.defer();
  //
  //   var filters = {
  //     label: ['com.docker.compose.project=' + name]
  //   };
  //
  //   ContainerService.containers(1, filters)
  //   .then(function success(data) {
  //     var containers = data;
  //     var services = StackHelper.getComposeV2ServicesFromContainers(containers);
  //     var stack = new StackV2ViewModel(name, services, containers);
  //     deferred.resolve(stack);
  //   })
  //   .catch(function error(err) {
  //     deferred.reject({ msg: 'Unable to retrieve stack details', err: err });
  //   });

    // return deferred.promise;
  // };

  service.stackV3 = function(name) {
    var deferred = $q.defer();

    var filters = {
      label: ['com.docker.stack.namespace=' + name]
    };

    $q.all({
      services: ServiceService.services(filters)
    })
    .then(function success(data) {
      var services = data.services;
      var stack = new StackV3ViewModel(name, services);
      deferred.resolve(stack);
    })
    .catch(function error(err) {
      deferred.reject({ msg: 'Unable to retrieve stack details', err: err });
    });

    return deferred.promise;
  };

  service.deleteStack = function(id) {
    return Stack.remove({ id: id }).$promise;
  };

  service.stackOperationUp = function(id) {
    return Stack.up({ id: id }).$promise;
  };

  service.stackOperationDown = function(id) {
    return Stack.down({ id: id }).$promise;
  };

  service.scaleService = function(id, serviceName, scale) {
    var payload = {
      ServiceName: serviceName,
      Scale: scale
    };
    return Stack.scale({ id: id }, payload).$promise;
  };

  return service;
}]);
