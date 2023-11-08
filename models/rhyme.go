package models

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"
)

type RhymeModel struct{}

// Rhyme represents a single rhyme. Rank allows the consumer
// to decide whether to show it or not.
type Rhyme struct {
	Word string `json:"word"`
	Rank int    `json:"rank"`
}

// GetRhymes gets the available rhymes for the passed word DB record
func (rm RhymeModel) GetRhymes(db *gorm.DB, word Word) ([]Rhyme, error) {

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	query, bindVals := buildRhymesSearchQuery(word)

	rhymes := make([]Rhyme, 0)


	rows, err := db.Raw(query, bindVals...).WithContext(ctx).Rows()
	if err != nil {
		log.Fatal(err)
		return rhymes, err
	}
	defer rows.Close()

	for rows.Next() {
		var rhyme Rhyme
		if err := rows.Scan(&rhyme.Word, &rhyme.Rank); err != nil {
			log.Fatal(err)
			return rhymes, err
		}
		rhymes = append(rhymes, rhyme)
	}

	return rhymes, nil
}
