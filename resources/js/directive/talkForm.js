/**
 * Created by languid on 4/2/15.
 */

define([
'jquery',
'angular',
'./module'
],
function($, ng, directive){
    directive
        .directive('talkForm', function(){
            return {
                link: function( scope, elem, attrs ){
                    var type = attrs.talkForm;
                    elem.find('button').click(function(){
                       scope.sendMessage(type)
                    });
                }
            }
        })
});