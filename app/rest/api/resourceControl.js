angular.module('dockm.rest')
.factory('ResourceControl', ['$resource', 'API_ENDPOINT_RESOURCE_CONTROLS', function ResourceControlFactory($resource, API_ENDPOINT_RESOURCE_CONTROLS) {
  'use strict';
  return $resource(API_ENDPOINT_RESOURCE_CONTROLS + '/:id', {}, {
    create: { method: 'POST' },
    get: { method: 'GET', params: { id: '@id' } },
    update: { method: 'PUT', params: { id: '@id' } },
    remove: { method: 'DELETE', params: { id: '@id'} }
  });
}]);
