package bootstrap

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"chatroom/config"
)

var methods = []func(){}

func Register(fn func()) {
	methods = append(methods, fn)
}

func Start(port string, onStart func()) {

	m := martini.Classic()

	m.Use(martini.Static("public"))
	m.Use(render.Renderer(render.Options{
		Charset: 	"UTF-8", // Sets encoding for json and html content-types. Default is "UTF-8".
		Delims:  	render.Delims{"${", "}"},
		Directory:	"resources/views",
	}))

	config.MappingController(m)

	http.Handle("/", m)

	onStart()

	for _, fn := range methods {
		go fn()
	}

	http.ListenAndServe(":"+port, nil)
}
