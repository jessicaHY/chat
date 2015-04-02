package bootstrap

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"chatroom/config"
	"runtime"
	"path"
)

var methods = []func(){}

func Register(fn func()) {
	methods = append(methods, fn)
}

func Start(port string, onStart func()) {

	resDir := "resources/views"
	_, filename, _, _ := runtime.Caller(1)

	m := martini.Classic()
	m.Use(martini.Static("public"))
	m.Use(martini.Static("assets"))
	m.Use(render.Renderer(render.Options{
		Charset: 	"UTF-8",
		Delims:  	render.Delims{"${", "}"},
		Directory:	path.Join(path.Dir(filename), resDir),
	}))

	config.MappingController(m)

	http.Handle("/", m)

	onStart()

	for _, fn := range methods {
		go fn()
	}

	http.ListenAndServe(":"+port, nil)
}
