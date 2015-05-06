/**
 * Created by languid on 4/2/15.
 */

requirejs.config({
    paths: {
        jquery: '/components/jquery/dist/jquery.min',
        kernel: 'core/kernel',
        react: '/components/react/react-with-addons.min',
        Backbone: '/components/backbone/backbone',
        underscore: '/components/underscore/underscore-min',
        Mustache : '/components/mustache/mustache.min'
    },
    shim: {
        Backbone: {
            deps: ['underscore']
        }
    }
});

require([
    'init'
]);