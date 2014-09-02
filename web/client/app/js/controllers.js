'use strict';

/* Controllers */

angular.module('godiwi.controllers', []).
controller('GodiController', ['apiURL', '$scope', '$location', '$resource', 'clientID',
    function NewGodiController(apiURL, $scope, $location, $resource, clientID) {
        var header = {"Client-ID": clientID};
        var State = $resource(apiURL + 'state', null, {
            defaults: {
                method: "DEFAULTS",
                headers: header,
            },
            update: {
                method: "PUT",
                headers: header,
            },
            get: {
                method: "GET",
                headers: header,
            },
            post: {
                method: "POST",
                headers: header,
            },
            "delete": {
                method: "DELETE",
                headers: header,
            },
        });

        var updateReadOnly = function(header) {
            if (!header("x-is-rw")){
                return;
            }
            $scope.stateReadOnly = header("x-is-rw") != 'true';
        };

        // These variables are kind of competing with each other if there are multple requests at once
        var updateDone = function(_, header) {
            $scope.isUpdating = false;
            $scope.updateFailed = false;
            if (header) {
                updateReadOnly(header);
            }
        };
        updateDone(); // init variables

        $scope.stateReadOnly = true;
        $scope.isUpdating = true; // we are updating now

        var updateFailed = function() {
            $scope.updateFailed = true;
        };

        // Will load up the websocket once we know the address, the first time we receive the state
        this.gotUpdateCallbackAt = Date.now();
        var parent = this;
        var firstStateHandler = function(state, header) {
            updateDone(null, header);

            if (!$scope.hasOwnProperty("$socket") || $scope.$socket.readyState != 1) {
                var conn = new WebSocket("ws://" + $location.host() + ':' + $location.port() + state.socketURL);
                conn.onmessage = function(val) {
                    var d = angular.fromJson(val.data);
                    if (d) {
                        if (d.state === 0 && d.clientID != clientID) { // state change
                            $scope.state.$get(null, function(val, header) {
                                console.log("WS FETCHED", clientID, d.clientID);
                                parent.gotUpdateCallbackAt = Date.now();
                                updateReadOnly(header);
                            });
                        }
                    }
                };
                // keep it around
                $scope.$socket = conn;
            }
        };

        $scope.default = State.defaults({}, updateDone, updateFailed);
        $scope.state = State.get({}, firstStateHandler, updateFailed);

        // Automatically put all changes, right when they happen
        $scope.$watch('state', function(nval) {
            if (parent.gotUpdateCallbackAt > Date.now() - 1000) {
                parent.gotUpdateCallbackAt = 0;
                return;
            }
            nval.$update({}, updateDone, updateFailed);
        }, true);

        return this;
    }
]).
controller("FilterController", ["$scope",
    function FilterController($scope) {
        // tracks whether a default is selected, one per default
        // We need this info to be part of the model, otherwise updates are difficult to handle
        $scope.fepDefaultSelections = [];
        var ctrl = this;

        // Called when a checkbox changes
        ctrl.onchange = function changed(name, add) {
            // We may be depending on the initial value coming from godi, lets be sure it has a value
            if (!$scope.state.fep) {
                $scope.state.fep = [];
            }
            if (add) {
                $scope.state.fep.push(name);
            } else {
                var idx = $scope.state.fep.indexOf(name);
                if (idx > -1) {
                    $scope.state.fep.splice(idx, 1);
                }
            }
        };

        
        ctrl.isSelected = function isSelected(filter) {
            if (!$scope.state.fep) {
                return false;
            }
            return $scope.state.fep.indexOf(filter) > -1;
        };

        ctrl.isDefault = function isDefault(filter) {
            if (!$scope.default.feps) {
                return false;
            }
            return $scope.default.feps.indexOf(filter) > -1;
        };

        ctrl.isNoDefault = function isNoDefault(filter) {
            return !ctrl.isDefault(filter);
        };

        ctrl.replace = function replace(index, nval) {
            $scope.state.fep[index] = nval;
        };

        ctrl.keyHandler = function keyPress(event, index, nval) {
            if (event.keyIdentifier == "Enter") {
                ctrl.replace(index, nval);
            }
        };

        $scope.$watch('state.fep', function fepChanged(nval) {
            if (!$scope.default.feps) {
                return;
            }
            for (var i = 0; i < $scope.default.feps.length; i++) {
                $scope.fepDefaultSelections[i] = ctrl.isSelected($scope.default.feps[i]);
            }
        }, true);
    }
]).
controller("LocationController", ["apiURL", "$scope", "$http", function LocationController(apiURL, $scope, $http) {
    var ctrl = this
    $scope.isValid = false
    this.listLocations = function listLocations(path) {
        return $http.get(apiURL + 'dirlist', {
            params: {
                'type': ($scope.type == 'source' && ($scope.mode == 'verify' ? 'sealOnly' : 'all')) || 'all', 
                'path': path,
            }}
        ).then(function success(res) {
            // It's never valid unless the item is actually selected
            $scope.isValid = false;
            return res.data;
        }, function error (code) {
            $scope.isValid = false;
            return []
        })
    }

    this.onSelect = function onSelect($item, $model, $label) {
        $scope.isValid = $scope.mode == 'verify' ? !$item.isDir : true
        $scope.newPath = $item.path
    }

    return this
}]);
