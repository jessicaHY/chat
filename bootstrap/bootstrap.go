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


	_, filename, _, _ := runtime.Caller(1)
	exeDir := path.Dir(filename)

	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Charset: 	"UTF-8",
		Delims:  	render.Delims{"${", "}"},
		Directory:	path.Join(exeDir, "resources/views"),
	}))

	m.Use(martini.Static(path.Join(exeDir, "public")))
	config.MappingController(m)

	http.Handle("/", m)

	onStart()

	for _, fn := range methods {
		go fn()
	}

	http.ListenAndServe(":"+port, nil)
}
