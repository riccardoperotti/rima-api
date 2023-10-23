package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SynonymsController holds methods and dependencies for the Synonyms Controller
type SynonymsController struct{}

func (s SynonymsController) GetSynonyms(c *gin.Context) {
	w := c.Param("word")

	ret := struct{ Message string }{"This will show synonyms for the word" + w}

	c.IndentedJSON(http.StatusOK, ret)
}
