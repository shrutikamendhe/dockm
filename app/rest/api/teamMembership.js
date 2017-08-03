angular.module('dockm.rest')
.factory('TeamMemberships', ['$resource', 'API_ENDPOINT_TEAM_MEMBERSHIPS', function TeamMembershipsFactory($resource, API_ENDPOINT_TEAM_MEMBERSHIPS) {
  'use strict';
  return $resource(API_ENDPOINT_TEAM_MEMBERSHIPS + '/:id/:action', {}, {
    create: { method: 'POST' },
    query: { method: 'GET', isArray: true },
    update: { method: 'PUT', params: { id: '@id' } },
    remove: { method: 'DELETE', params: { id: '@id'} }
  });
}]);
