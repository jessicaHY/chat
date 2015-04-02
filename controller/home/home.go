package home

import (
	_ "github.com/antonholmquist/jason"
	"github.com/martini-contrib/render"
)

func Home(rend render.Render){
	rend.HTML(200, "index", nil)
}
