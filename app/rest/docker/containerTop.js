angular.module('dockm.rest')
.factory('ContainerTop', ['$http', 'API_ENDPOINT_ENDPOINTS', 'EndpointProvider', function ($http, API_ENDPOINT_ENDPOINTS, EndpointProvider) {
  'use strict';
  return {
    get: function (id, params, callback, errorCallback) {
      $http({
        method: 'GET',
        url: API_ENDPOINT_ENDPOINTS + '/' + EndpointProvider.endpointID() + '/docker/containers/' + id + '/top',
        params: {
          ps_args: params.ps_args
        }
      }).success(callback).error(function (data, status, headers, config) {
        console.log(data);
      });
    }
  };
}]);
