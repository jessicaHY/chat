/**
 * Created by Yinxiong on 2014-11-05.
 */

define([
'jquery',
'kernel',
'Backbone'
],
function( $, core, Backbone ){

    var Socket = Backbone.Model.extend({
        isListen: false,
        isConnected: false,
        _retryTimes: 0,
        defaults: {
            reconnect: true,
            arg: '',
            events: {
                message: {},
                open: $.noop,
                error: $.noop,
                close: $.noop,
                giveup: $.noop
            }
        },
        socket: null,
        listen: function(){
            var self = this;
            if( this.isListen ){
                return
            }
            this.isListen = true;
            this.socket = new WebSocket('ws://'+location.host+'/socket/'+this.get('arg'));

            this.socket.onmessage = function( message ){
                var data = JSON.parse(message.data);
                if( data.method in self.get('events').message ){
                    self.get('events').message[data.method]( data.data, message )
                }
            };

            this.socket.onopen = function(e){
                self.get('events').open.call(self, e);
                self._retryTimes = 0;
                self.isConnected = true;
            };

            this.socket.onclose = function(e){
                self.get('events').close.call(self, e);
                self.get('reconnect') && self.reconnect();
                self.isConnected = false;
            }
        },
        reconnect: function(){
            if( this._retryTimes >= 3 ){
                this.get('events').giveup();
                this.disconnect(true);
                this.isListen = false;
                return;
            }
            this._retryTimes++;
            console.warn('websocket disconnect, retry', this._retryTimes);
            this.socket.close();
            this.isListen = false;
            this.listen();
        },
        disconnect: function( force ){
            if( force ){
                this.set('reconnect', false);
            }
            this.socket.close( force );
        },
        emit: function( method, data ){
            return this.send({
                method: method,
                data: data
            })
        },
        broadcast: function( method, data ){
            return this.emit('broadcast', {
                method: method,
                data: data
            });
        },
        on: function( method, fn ){
            if( typeof method != 'string' ) return;

            if( (/^(error|close|open|giveup)$/).test(method) ){
                this.get('events')[method] = fn;
            }else{
                this.get('events').message[method] = fn;
            }
        },
        send: function( data ){
            return this.socket.send(data)
        }
    });

    return function( events, reconnect ){
        return new Socket( events, reconnect )
    }
});