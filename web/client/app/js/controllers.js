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
controller('GodiController', ['$scope', '$resource',
    function NewGodiController($scope, $resource) {
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


        $scope.state = State.get({}, updateDone, updateFailed);
        $scope.default = State.defaults({}, updateDone, updateFailed);

        // Set suitable defaults

        // Automatically put all changes, right when they happen
        $scope.$watchCollection('state', function(nval, oval) {
        	nval.$update({}, updateDone, updateFailed);
        });

        return this;
    }
]);
