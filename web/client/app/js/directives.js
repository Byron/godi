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
}).
directive('pathinput', function() {
    return {
        restrict: 'E',
        templateUrl: 'pathTemplate.html',
        scope: {
            mode: '=',
            paths: '=',
        },
        link:  function link(scope, element, attrs) {
            scope.title = attrs['title'] || "TITLE UNSET"
            scope.type = attrs['type'] || "source"
        }
    }
}).
directive("keepScroll",  function(){
  return {
    controller : function($scope){
      var element = 0;
      
      this.setElement = function(el){
        element = el;
      }

      this.addItem = function(item){
        element.scrollTop = item.offsetTop; //1px for margin
      };
    },
    
    link : function(scope,el,attr, ctrl) {
        ctrl.setElement(el[0]);
    }
  };
}).
directive("scrollItem", function(){
  return{
    require : "^keepScroll",
    link : function(scope, el, att, scrCtrl){
        scrCtrl.addItem(el[0]);
    }
  }
});
