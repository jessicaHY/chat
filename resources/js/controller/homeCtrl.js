/**
 * Created by languid on 4/2/15.
 */
define([
'kernel',
'jquery',
'angular',
'service/SocketInstance',
'directive/talkForm',
'ngSanitize'
],
function(core, $, ng){

    var App = ng.module('App', ['App.services', 'App.directives', 'ngSanitize']);

    App
        .value('UserType', {
            Author: 1,
            User: 2
        })
        .controller('homeCtrl', ['$scope','SocketInstance', 'UserType', 'Helper', function($scope, SocketInstance, UserType, Helper){

            //exports
            $scope.UserType = UserType;

            //set model
            $scope.authorContent = '';
            $scope.userContent = '';

            $scope.authorTalkList = [];
            $scope.userTalkList = [];

            SocketInstance.setScope( $scope );
            SocketInstance.on('authorMessage', function( data ){
                $scope.authorTalkList.push(data);
            });
            SocketInstance.on('userMessage', function( data ){
                $scope.userTalkList.push(data);
            });

            $scope.sendMessage = function( type ){
                if( type == UserType.Author ){
                    return SocketInstance.emit('authorSend', {
                        content: $scope.authorContent
                    });
                } else if ( type == UserType.User ){
                    return SocketInstance.emit('userSend', {
                        content: $scope.userContent
                    });
                }
            }

        }]);

    ng.bootstrap(document, ['App']);
});