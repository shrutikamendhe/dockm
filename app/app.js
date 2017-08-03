angular.module('dockm.filters', []);
angular.module('dockm.rest', ['ngResource']);
angular.module('dockm.services', []);
angular.module('dockm.helpers', []);
angular.module('dockm', [
  'ui.bootstrap',
  'ui.router',
  'isteven-multi-select',
  'ngCookies',
  'ngSanitize',
  'ngFileUpload',
  'angularUtils.directives.dirPagination',
  'LocalStorageModule',
  'angular-jwt',
  'angular-google-analytics',
  'dockm.templates',
  'dockm.filters',
  'dockm.rest',
  'dockm.helpers',
  'dockm.services',
  'auth',
  'dashboard',
  'container',
  'containerConsole',
  'containerLogs',
  'serviceLogs',
  'containers',
  'createContainer',
  'createNetwork',
  'createRegistry',
  'createSecret',
  'createService',
  'createStack',
  'createVolume',
  'docker',
  'endpoint',
  'endpointAccess',
  'endpointInit',
  'endpoints',
  'events',
  'image',
  'images',
  'main',
  'network',
  'networks',
  'node',
  'registries',
  'registry',
  'registryAccess',
  'secrets',
  'secret',
  'service',
  'services',
  'unregisteredstackv2',
  'stackv2',
  'stackv3',
  'stacks',
  'settings',
  'sidebar',
  'stats',
  'swarm',
  'task',
  'team',
  'teams',
  'templates',
  'user',
  'users',
  'userSettings',
  'volume',
  'volumes'])
  .config(['$stateProvider', '$urlRouterProvider', '$httpProvider', 'localStorageServiceProvider', 'jwtOptionsProvider', 'AnalyticsProvider', '$uibTooltipProvider', '$compileProvider', function ($stateProvider, $urlRouterProvider, $httpProvider, localStorageServiceProvider, jwtOptionsProvider, AnalyticsProvider, $uibTooltipProvider, $compileProvider) {
    'use strict';

    var environment = '@@ENVIRONMENT';
    if (environment === 'production') {
      $compileProvider.debugInfoEnabled(false);
    }

    localStorageServiceProvider
    .setPrefix('dockm');

    jwtOptionsProvider.config({
      tokenGetter: ['LocalStorage', function(LocalStorage) {
        return LocalStorage.getJWT();
      }],
      unauthenticatedRedirector: ['$state', function($state) {
        $state.go('auth', {error: 'Your session has expired'});
      }]
    });
    $httpProvider.interceptors.push('jwtInterceptor');

    AnalyticsProvider.setAccount('@@CONFIG_GA_ID');
    AnalyticsProvider.startOffline(true);

    $urlRouterProvider.otherwise('/auth');

    toastr.options.timeOut = 3000;

    $uibTooltipProvider.setTriggers({
      'mouseenter': 'mouseleave',
      'click': 'click',
      'focus': 'blur',
      'outsideClick': 'outsideClick'
    });

    $stateProvider
    .state('root', {
      abstract: true,
      resolve: {
        requiresLogin: ['StateManager', function (StateManager) {
          var applicationState = StateManager.getState();
          return applicationState.application.authentication;
        }]
      }
    })
    .state('auth', {
      parent: 'root',
      url: '/auth',
      params: {
        logout: false,
        error: ''
      },
      views: {
        'content@': {
          templateUrl: 'app/components/auth/auth.html',
          controller: 'AuthenticationController'
        }
      },
      data: {
        requiresLogin: false
      }
    })
    .state('containers', {
      parent: 'root',
      url: '/containers/',
      views: {
        'content@': {
          templateUrl: 'app/components/containers/containers.html',
          controller: 'ContainersController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('container', {
      url: '^/containers/:id',
      views: {
        'content@': {
          templateUrl: 'app/components/container/container.html',
          controller: 'ContainerController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('stats', {
      url: '^/containers/:id/stats',
      views: {
        'content@': {
          templateUrl: 'app/components/stats/stats.html',
          controller: 'StatsController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('containerlogs', {
      url: '^/containers/:id/logs',
      views: {
        'content@': {
          templateUrl: 'app/components/containerLogs/containerlogs.html',
          controller: 'ContainerLogsController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('servicelogs', {
      url: '^/services/:id/logs',
      views: {
        'content@': {
          templateUrl: 'app/components/serviceLogs/servicelogs.html',
          controller: 'ServiceLogsController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('console', {
      url: '^/containers/:id/console',
      views: {
        'content@': {
          templateUrl: 'app/components/containerConsole/containerConsole.html',
          controller: 'ContainerConsoleController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('dashboard', {
      parent: 'root',
      url: '/dashboard',
      views: {
        'content@': {
          templateUrl: 'app/components/dashboard/dashboard.html',
          controller: 'DashboardController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('actions', {
      abstract: true,
      url: '/actions',
      views: {
        'content@': {
          template: '<div ui-view="content@"></div>'
        },
        'sidebar@': {
          template: '<div ui-view="sidebar@"></div>'
        }
      }
    })
    .state('actions.create', {
      abstract: true,
      url: '/create',
      views: {
        'content@': {
          template: '<div ui-view="content@"></div>'
        },
        'sidebar@': {
          template: '<div ui-view="sidebar@"></div>'
        }
      }
    })
    .state('actions.create.container', {
      url: '/container',
      views: {
        'content@': {
          templateUrl: 'app/components/createContainer/createcontainer.html',
          controller: 'CreateContainerController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('actions.create.network', {
      url: '/network',
      views: {
        'content@': {
          templateUrl: 'app/components/createNetwork/createnetwork.html',
          controller: 'CreateNetworkController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('actions.create.registry', {
      url: '/registry',
      views: {
        'content@': {
          templateUrl: 'app/components/createRegistry/createregistry.html',
          controller: 'CreateRegistryController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('actions.create.secret', {
      url: '/secret',
      views: {
        'content@': {
          templateUrl: 'app/components/createSecret/createsecret.html',
          controller: 'CreateSecretController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('actions.create.service', {
      url: '/service',
      views: {
        'content@': {
          templateUrl: 'app/components/createService/createservice.html',
          controller: 'CreateServiceController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('actions.create.stack', {
      url: '/stack',
      views: {
        'content@': {
          templateUrl: 'app/components/createStack/createstack.html',
          controller: 'CreateStackController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('actions.create.volume', {
      url: '/volume',
      views: {
        'content@': {
          templateUrl: 'app/components/createVolume/createvolume.html',
          controller: 'CreateVolumeController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('docker', {
      url: '/docker/',
      views: {
        'content@': {
          templateUrl: 'app/components/docker/docker.html',
          controller: 'DockerController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('endpoints', {
      url: '/endpoints/',
      views: {
        'content@': {
          templateUrl: 'app/components/endpoints/endpoints.html',
          controller: 'EndpointsController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('endpoint', {
      url: '^/endpoints/:id',
      views: {
        'content@': {
          templateUrl: 'app/components/endpoint/endpoint.html',
          controller: 'EndpointController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('endpoint.access', {
      url: '^/endpoints/:id/access',
      views: {
        'content@': {
          templateUrl: 'app/components/endpointAccess/endpointAccess.html',
          controller: 'EndpointAccessController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('endpointInit', {
      url: '/init/endpoint',
      views: {
        'content@': {
          templateUrl: 'app/components/endpointInit/endpointInit.html',
          controller: 'EndpointInitController'
        }
      }
    })
    .state('events', {
      url: '/events/',
      views: {
        'content@': {
          templateUrl: 'app/components/events/events.html',
          controller: 'EventsController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('images', {
      url: '/images/',
      views: {
        'content@': {
          templateUrl: 'app/components/images/images.html',
          controller: 'ImagesController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('image', {
      url: '^/images/:id/',
      views: {
        'content@': {
          templateUrl: 'app/components/image/image.html',
          controller: 'ImageController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('networks', {
      url: '/networks/',
      views: {
        'content@': {
          templateUrl: 'app/components/networks/networks.html',
          controller: 'NetworksController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('network', {
      url: '^/networks/:id/',
      views: {
        'content@': {
          templateUrl: 'app/components/network/network.html',
          controller: 'NetworkController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('node', {
      url: '^/nodes/:id/',
      views: {
        'content@': {
          templateUrl: 'app/components/node/node.html',
          controller: 'NodeController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('registries', {
      url: '/registries/',
      views: {
        'content@': {
          templateUrl: 'app/components/registries/registries.html',
          controller: 'RegistriesController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('registry', {
      url: '^/registries/:id',
      views: {
        'content@': {
          templateUrl: 'app/components/registry/registry.html',
          controller: 'RegistryController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('registry.access', {
      url: '^/registries/:id/access',
      views: {
        'content@': {
          templateUrl: 'app/components/registryAccess/registryAccess.html',
          controller: 'RegistryAccessController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('secrets', {
      url: '^/secrets/',
      views: {
        'content@': {
          templateUrl: 'app/components/secrets/secrets.html',
          controller: 'SecretsController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('secret', {
      url: '^/secret/:id/',
      views: {
        'content@': {
          templateUrl: 'app/components/secret/secret.html',
          controller: 'SecretController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('services', {
      url: '/services/',
      views: {
        'content@': {
          templateUrl: 'app/components/services/services.html',
          controller: 'ServicesController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('service', {
      url: '^/service/:id/',
      views: {
        'content@': {
          templateUrl: 'app/components/service/service.html',
          controller: 'ServiceController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('stacks', {
      url: '/stacks/',
      views: {
        'content@': {
          templateUrl: 'app/components/stacks/stacks.html',
          controller: 'StacksController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('stack', {
      abstract: true,
      url: '/stacks',
      views: {
        'content@': {
          template: '<div ui-view="content@"></div>'
        },
        'sidebar@': {
          template: '<div ui-view="sidebar@"></div>'
        }
      }
    })
    .state('stack.v2', {
      url: '^/stacks/v2/:id/',
      views: {
        'content@': {
          templateUrl: 'app/components/stack/v2/stack.html',
          controller: 'StackV2Controller'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('stack.v3', {
      url: '^/stacks/v3/:id/',
      views: {
        'content@': {
          templateUrl: 'app/components/stack/v3/stack.html',
          controller: 'StackV3Controller'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('stack.v2.unregistered', {
      url: '^/stacks/v2/unregistered/:name/',
      views: {
        'content@': {
          templateUrl: 'app/components/stack/unregisteredv2/stack.html',
          controller: 'UnregisteredStackV2Controller'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('settings', {
      url: '/settings/',
      views: {
        'content@': {
          templateUrl: 'app/components/settings/settings.html',
          controller: 'SettingsController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('task', {
      url: '^/task/:id',
      views: {
        'content@': {
          templateUrl: 'app/components/task/task.html',
          controller: 'TaskController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('templates', {
      url: '/templates/',
      params: {
        key: 'containers',
        hide_descriptions: false
      },
      views: {
        'content@': {
          templateUrl: 'app/components/templates/templates.html',
          controller: 'TemplatesController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('templates_linuxserver', {
      url: '^/templates/linuxserver.io',
      params: {
        key: 'linuxserver.io',
        hide_descriptions: true
      },
      views: {
        'content@': {
          templateUrl: 'app/components/templates/templates.html',
          controller: 'TemplatesController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('volumes', {
      url: '/volumes/',
      views: {
        'content@': {
          templateUrl: 'app/components/volumes/volumes.html',
          controller: 'VolumesController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('volume', {
      url: '^/volumes/:id',
      views: {
        'content@': {
          templateUrl: 'app/components/volume/volume.html',
          controller: 'VolumeController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('users', {
      url: '/users/',
      views: {
        'content@': {
          templateUrl: 'app/components/users/users.html',
          controller: 'UsersController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('user', {
      url: '^/users/:id',
      views: {
        'content@': {
          templateUrl: 'app/components/user/user.html',
          controller: 'UserController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('userSettings', {
      url: '/userSettings/',
      views: {
        'content@': {
          templateUrl: 'app/components/userSettings/userSettings.html',
          controller: 'UserSettingsController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('teams', {
      url: '/teams/',
      views: {
        'content@': {
          templateUrl: 'app/components/teams/teams.html',
          controller: 'TeamsController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('team', {
      url: '^/teams/:id',
      views: {
        'content@': {
          templateUrl: 'app/components/team/team.html',
          controller: 'TeamController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    })
    .state('swarm', {
      url: '/swarm/',
      views: {
        'content@': {
          templateUrl: 'app/components/swarm/swarm.html',
          controller: 'SwarmController'
        },
        'sidebar@': {
          templateUrl: 'app/components/sidebar/sidebar.html',
          controller: 'SidebarController'
        }
      }
    });
  }])
  .run(['$rootScope', '$state', 'Authentication', 'authManager', 'StateManager', 'EndpointProvider', 'Notifications', 'Analytics', function ($rootScope, $state, Authentication, authManager, StateManager, EndpointProvider, Notifications, Analytics) {
    EndpointProvider.initialize();
    StateManager.initialize().then(function success(state) {
      if (state.application.authentication) {
        authManager.checkAuthOnRefresh();
        authManager.redirectWhenUnauthenticated();
        Authentication.init();
        $rootScope.$on('tokenHasExpired', function($state) {
          $state.go('auth', {error: 'Your session has expired'});
        });
      }
      if (state.application.analytics) {
        Analytics.offline(false);
        Analytics.registerScriptTags();
        Analytics.registerTrackers();
        $rootScope.$on('$stateChangeSuccess', function (event, toState, toParams, fromState, fromParams) {
          Analytics.trackPage(toState.url);
          Analytics.pageView();
        });
      }
    }, function error(err) {
      Notifications.error('Failure', err, 'Unable to retrieve application settings');
    });

    $rootScope.$state = $state;
  }])
  // This is your docker url that the api will use to make requests
  // You need to set this to the api endpoint without the port i.e. http://192.168.1.9
  // .constant('DOCKER_PORT', '') // Docker port, leave as an empty string if no port is required.  If you have a port, prefix it with a ':' i.e. :4243
  .constant('API_ENDPOINT_AUTH', 'api/auth')
  .constant('API_ENDPOINT_DOCKERHUB', 'api/dockerhub')
  .constant('API_ENDPOINT_ENDPOINTS', 'api/endpoints')
  .constant('API_ENDPOINT_REGISTRIES', 'api/registries')
  .constant('API_ENDPOINT_RESOURCE_CONTROLS', 'api/resource_controls')
  .constant('API_ENDPOINT_SETTINGS', 'api/settings')
  .constant('API_ENDPOINT_STATUS', 'api/status')
  .constant('API_ENDPOINT_USERS', 'api/users')
  .constant('API_ENDPOINT_TEAMS', 'api/teams')
  .constant('API_ENDPOINT_TEAM_MEMBERSHIPS', 'api/team_memberships')
  .constant('API_ENDPOINT_TEMPLATES', 'api/templates')
  .constant('DEFAULT_TEMPLATES_URL', 'https://raw.githubusercontent.com/Click2Cloud/templates/master/templates.json')
  .constant('PAGINATION_MAX_ITEMS', 10);
