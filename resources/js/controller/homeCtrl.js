/**
 * Created by languid on 4/2/15.
 */
define([
'kernel',
'jquery',
'underscore',
'Backbone',
'Mustache',
'service/Websocket'
],
function(core, $, _, Backbone, Mustache, Websocket){

    var User_Type = {
        Author: 1,
        User: 2
    };

    var Content_Type = {
        Chapter: 0,
        Reply: 1
    };

    var socket = new Websocket({
        arg: '1'
    });


    var TalkItemModel = Backbone.Model.extend({
        defaults: {
            user_type: User_Type.User,
            content: '',
            createTime: '',
            id: 0,
            type: 0,
            userInfo: {
                icon: '',
                id: 0,
                isAuthor: false,
                name: '',
                subscribed: false
            }
        }
    });

    var TalkCollection = Backbone.Collection.extend({
        model: TalkItemModel
    });

    var TalkItemView = Backbone.View.extend({
        tagName: 'LI',
        template: function(){
            var authorTpl = $('#Template_AuthorTalkItem').html();
            var userTpl = $('#Template_UserTalkItem').html();
            return function( t ){
                return t == User_Type.Author ? authorTpl : userTpl;
            }
        }(),
        render: function(){
            var m = this.model.toJSON();
            var html = Mustache.render(
                this.template( this.model.get('user_type') ), m
            );
            this.$el.html(html);
            return this;
        }
    });

    var TalkListView = Backbone.View.extend({
        type: 0,
        initialize: function( opt ){
            var self = this;
            this.type = opt.type;
            this.collection.bind('add', function( model ){
                model.set('user_type', self.type);
                var view = new TalkItemView({ model: model });
                self.$el.append( view.render().$el );
            })
        }
    });

    var FormView = Backbone.View.extend({
        events: {
            'keyup textarea': 'check',
            'click button': 'send'
        },
        type: '',
        value: '',
        initialize: function( opts ){
            this.type = opts.type;
            this.$textarea = this.$el.find('textarea');
            this.$btn = this.$el.find('.btn');
        },
        check: function( e ){
            var value = this.$textarea.val().trim();
            if( value ){
                this.$btn.removeClass('disabled')
            } else {
                this.$btn.addClass('disabled')
            }
            this.value = value;
        },
        send: function(){
            if( this.type == User_Type.Author ){
                socket.emit('authorSend', {
                    content: this.value
                });
            } else if( this.type == User_Type.User ){
                socket.emit('userSend', {
                    content: this.value
                });
            }
        }
    });

    var App = Backbone.View.extend({
        el: '#App',
        initialize: function(){
            var self = this;
            this.author_form = new FormView({
                type: User_Type.Author,
                el: this.$el.find('#AuthorForm')
            });
            this.user_form = new FormView({
                type: User_Type.User,
                el: this.$el.find('#UserForm')
            });

            this.author_list = new TalkListView({
                type: User_Type.Author,
                collection: new TalkCollection(),
                el: this.$el.find('#AuthorTalk ul')
            });

            this.user_list = new TalkListView({
                type: User_Type.User,
                collection: new TalkCollection(),
                el: this.$el.find('#UserTalk ul')
            });

            socket.on('authorMessage', function( data ){
                self.author_list.collection.add(data);
            });
            socket.on('userMessage', function( data ){
                self.user_list.collection.add(data);
            });
        }
    });

    socket.listen();

    return new App;

});