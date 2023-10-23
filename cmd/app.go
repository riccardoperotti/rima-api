package main

import (
	"github.com/gin-gonic/gin"
	"github.com/riccardoperotti/rima-api/controllers"
)

func main() {
	router := gin.Default()

	// TODO: fetch dependencies (db, config, etc) and add them to
	// each controller instance

	rhymes := new(controllers.RhymesController)
	router.GET("/rima/:word", rhymes.GetRhymes)

	synonyms := new(controllers.SynonymsController)
	router.GET("/sinonimo/:word", synonyms.GetSynonyms)

	router.SetTrustedProxies(nil)
	router.Run()
}
