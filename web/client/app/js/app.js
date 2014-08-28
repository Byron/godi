'use strict';


// Declare app level module which depends on filters, and services
// godiwi == godi-web-interface
angular.module('godiwi', [
  'ngRoute',
  'ngResource',
  'ui.bootstrap',
  'godiwi.filters',
  'godiwi.services',
  'godiwi.directives',
  'godiwi.controllers',
]).
config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/view1', {templateUrl: 'partials/partial1.html'});
  $routeProvider.when('/view2', {templateUrl: 'partials/partial2.html'});
  $routeProvider.otherwise({redirectTo: '/view1'});
}]);
