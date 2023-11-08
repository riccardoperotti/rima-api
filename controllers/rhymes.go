package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/riccardoperotti/rima-api/db"
	"github.com/riccardoperotti/rima-api/models"
	"gorm.io/gorm"
)

// RhymesController holds methods and dependencies for the Rhymes controller
type RhymesController struct{}

type Result struct {
	Lookup struct {
		Word      string `json:"word"`
		Type      string `json:"type"`
		Syllables string `json:"syllables"`
		Sounds    string `json:"sounds"`
	} `json:"lookup"`
	Count  int            `json:"count"`
	Error  string         `json:"error"`
	Rhymes []models.Rhyme `json:"rhymes"`
}

var TypeStrings = map[string]string{
	"A": "Aguda",
	"G": "Grave",
	"E": "Esdrújula",
	"S": "Sobreesdrújula",
}

func (r RhymesController) GetRhymes(c *gin.Context) {
	w := strings.ToLower(c.Param("word"))

	res := Result{}
	res.Lookup.Word = w

	// connect to the db
	dbh, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo establecer una conexión. Por favor inténtelo más tarde."})
		log.Printf("Error connecting to DB: %s", err)
		return
	}

	// find this word in the database
	var wordModel = new(models.WordModel)

	word, err := wordModel.GetWord(dbh, w)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			res.Error = fmt.Sprintf("La palabra '%s' no está en nuestra base de datos.", w)
			c.IndentedJSON(http.StatusOK, res)
			return
		}
		log.Fatal(err)
		res.Error = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res.Lookup.Syllables = strings.Join(word.Syllables(), "-")
	res.Lookup.Sounds = strings.Join(word.Sounds(), "-")
	res.Lookup.Type = TypeStrings[word.Type]

	// get rhymes for this word
	var rhymeModel = new(models.RhymeModel)

	rhymes, err := rhymeModel.GetRhymes(dbh, word)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			res.Error = fmt.Sprintf("No se encontraron rimas para la palabra '%s'.", w)
			c.IndentedJSON(http.StatusOK, res)
			return
		}

		log.Printf("Error fetching rhymes: %s.", err)
		res.Error = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res.Rhymes = rhymes
	res.Count = len(rhymes)

	c.IndentedJSON(http.StatusOK, res)
}
