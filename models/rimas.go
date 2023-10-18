package models

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type RimasModel struct{}

// Rima represents a single rhyme. Rank allows the consumer
// to decide whether to show it or not.
type Rima struct {
	Rima string `json:"rima"`
	Rank int    `json:"rank"`
}

// GetRimas gets the available rhymes for the passed word DB record
func (rm RimasModel) GetRimas(dbh *sql.DB, word Word) ([]Rima, error) {

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	query, bindVals := buildRimasSearchQuery(word)

	rimas := make([]Rima, 0)

	rows, err := dbh.QueryContext(ctx, query, bindVals...)
	if err != nil {
		log.Fatal(err)
		return rimas, err
	}
	defer rows.Close()

	for rows.Next() {
		var rima Rima
		if err := rows.Scan(&rima.Rima, &rima.Rank); err != nil {
			log.Fatal(err)
			return rimas, err
		}
		rimas = append(rimas, rima)
	}

	return rimas, nil
}
