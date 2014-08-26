'use strict';

/* Directives */


angular.module('godiwi.directives', []).
directive('appVersion', ['version',
    function(version) {
        return function(scope, elm, attrs) {
            elm.text(version);
        };
    }
]).
directive('uniqueAndNoFepDefault', function() {
    return {
        require: "ngModel",
        link: function(scope, elm, attrs, ctrl) {
            ctrl.$parsers.unshift(function(viewValue) {
                var isValid = false;
                if (viewValue) {
                    isValid = scope.default.feps.indexOf(viewValue.toUpperCase()) < 0;
                    if (isValid) {
                        // it's not a default value - check if it's already in the file exclude patterns list
                        for (var idx in scope.state.fep) {
                            if (scope.state.fep[idx].toLowerCase() == viewValue.toLowerCase()) {
                                isValid = false;
                                break;
                            }
                        }
                    }
                }
                ctrl.$setValidity('noFepDefault', isValid);
                if (isValid) {
                    return viewValue;
                } else {
                    return undefined;
                }
            });
        }
    };
});
