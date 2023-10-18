package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SinonimosController holds methods and dependencies for the Sinonimos controller
type SinonimosController struct{}

func (s SinonimosController) GetSinonimos(c *gin.Context) {
	w := c.Param("word")

	ret := struct{ Message string }{"This will show sinonimos for " + w}

	c.IndentedJSON(http.StatusOK, ret)
}
