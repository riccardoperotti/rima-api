package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/riccardoperotti/rima-api/db"
	"github.com/riccardoperotti/rima-api/models"
)

// RimasController holds methods and dependencies for the Rimas controller
type RimasController struct{}

type Result struct {
	Word struct {
		Name      string `json:"name"`
		Type      string `json:"type"`
		Syllables string `json:"syllables"`
		Sounds    string `json:"sounds"`
	} `json:"word"`
	Count int           `json:"count"`
	Error string        `json:"error"`
	Rimas []models.Rima `json:"rimas"`
}

func (r RimasController) GetRimas(c *gin.Context) {
	w := strings.ToLower(c.Param("word"))

	res := Result{}
	res.Word.Name = w

	// connect to the db
	dbh, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect. Please try again later."})
		log.Printf("Error connecting to DB: %s", err)
		return
	}
	defer dbh.Close()

	// find this word in the database
	var wordModel = new(models.WordModel)

	word, err := wordModel.GetWord(dbh, w)
	if err != nil {
		if err == sql.ErrNoRows {
			res.Error = fmt.Sprintf("La palabra '%s' no está en nuestra base de datos.", w)
			c.IndentedJSON(http.StatusOK, res)
			return
		}
		log.Fatal(err)
		res.Error = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	// TODO: add more info about the word (analytics, sounds, etc)
	res.Word.Syllables = strings.Join(word.Syllables(), "-")
	res.Word.Sounds = strings.Join(word.Sounds(), "-")
	res.Word.Type = word.Type

	// get rimas for this word
	var rimasModel = new(models.RimasModel)

	rimas, err := rimasModel.GetRimas(dbh, word)
	if err != nil {
		if err == sql.ErrNoRows {
			res.Error = fmt.Sprintf("No se encontraron rimas para la palabra '%s'.", w)
			c.IndentedJSON(http.StatusOK, res)
			return
		}

		log.Printf("Error fetching rimas: %s.", err)
		res.Error = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res.Rimas = rimas
	res.Count = len(rimas)

	c.IndentedJSON(http.StatusOK, res)
}
