'use strict';

/* Services */


// Demonstrate how to register services
// In this case it is a simple value service.
angular.module('godiwi.services', []).
  value('version', '0.2').
  constant('clientID', (Math.floor(Math.random() * 0x20000) + Date.now().toString(16)))
  