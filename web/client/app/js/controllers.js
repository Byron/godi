'use strict';

/* Controllers */

angular.module('godiwi.controllers', [])
    .controller('MyCtrl1', ['$scope',
        function($scope) {

        }
    ])
    .controller('MyCtrl2', ['$scope',
        function($scope) {

        }
    ]).
controller('GodiController', ['$scope', '$location', '$resource',
    function NewGodiController($scope, $location, $resource) {
        var State = $resource('/api/v1/state', null, {
            defaults: {
                method: "DEFAULTS"
            },
            update: {
                method: "PUT"
            }
        });

        // These variables are kind of competing with each other if there are multple requests at once
        var updateDone = function() {
            $scope.isUpdating = false;
            $scope.updateFailed = false;
        };
        updateDone(); // init variables
        $scope.isUpdating = true; // we are updating now

        var updateFailed = function() {
            $scope.updateFailed = true;
        };

        // Will load up the websocket once we know the address, the first time we receive the state
        var firstStateHandler = function(state) {
            updateDone();

            if (!$scope.hasOwnProperty("$socket") || $scope.$socket.readyState != 1) {
                var conn = new WebSocket("ws://" + $location.host() + ':' + $location.port() + state.socketURL);
                conn.onmessage = function() {
                	// This will also cause us to respond to events triggered by ourselves.
                	// Could be an optimization to prevent that ... lets see.
                    $scope.state.$get();
                };
                // keep it around
                $scope.$socket = conn;
            }
        };

        $scope.default = State.defaults({}, updateDone, updateFailed);
        $scope.state = State.get({}, firstStateHandler, updateFailed);

        // Automatically put all changes, right when they happen
        $scope.$watchCollection('state', function(nval, oval) {
            nval.$update({}, updateDone, updateFailed);
        });

        return this;
    }
]);
