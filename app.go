package main

import (
	"github.com/gin-gonic/gin"
	"github.com/riccardoperotti/rima-api/controllers"
)

func main() {
	router := gin.Default()

	// TODO: fetch dependencies (db, config, etc) and add them to
	// each controller instance

	rimas := new(controllers.RimasController)
	router.GET("/rima/:word", rimas.GetRimas)

	sinonimos := new(controllers.SinonimosController)
	router.GET("/sinonimo/:word", sinonimos.GetSinonimos)

	router.SetTrustedProxies(nil)
	router.Run()
}
