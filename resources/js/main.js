/**
 * Created by languid on 4/2/15.
 */

requirejs.config({
    paths: {
        jquery: '/components/jquery/dist/jquery.min',
        kernel: 'core/kernel',
        angular: '/components/angular/angular.min',
        react: '/components/react/react-with-addons.min',
        ngSanitize: '/components/angular-sanitize/angular-sanitize.min',
        ngWebSocket: '/components/angular-websocket/angular-websocket.min'
    },
    shim: {
        ngSanitize: {
            deps: ['angular']
        },
        ngWebSocket: {
            deps: ['angular']
        },
        angular: {
            exports: 'angular'
        }
    }
});

require([
    'init'
]);