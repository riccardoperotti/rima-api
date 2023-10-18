package models

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type WordModel struct{}

// Word represents the a row in the `lexico` table
type Word struct {
	Name          string
	SyllableCount int
	Type          string
	Syllable4     string
	Syllable3     string
	Syllable2     string
	Syllable1     string
	EndsWith      string
	Rank          int
}

// GetWord fetches this word's complete record from the `lexico` table
func (wm WordModel) GetWord(dbh *sql.DB, word string) (Word, error) {

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	q := `
		select	palabra,
				silabas,
				tipo,
				COALESCE(silaba4, ""),
				COALESCE(silaba3, ""),
				COALESCE(silaba2, ""),
				COALESCE(silaba1, ""),
				COALESCE(final, ""),
				rank
		from 	lexico
		where 	palabra = ?
	`

	stmt, err := dbh.PrepareContext(ctx, q)
	if err != nil {
		log.Printf("Error when preparing SQL statement: %s", err.Error())
		return Word{}, err
	}
	defer stmt.Close()

	var dbW Word
	row := stmt.QueryRowContext(ctx, word)
	err = row.Scan(
		&dbW.Name,
		&dbW.SyllableCount,
		&dbW.Type,
		&dbW.Syllable4,
		&dbW.Syllable3,
		&dbW.Syllable2,
		&dbW.Syllable1,
		&dbW.EndsWith,
		&dbW.Rank,
	)
	if err != nil {
		log.Printf("Error when scanning results:  %s", err.Error())
		return Word{}, err
	}

	return dbW, nil
}

// Syllables returns the non-empty syllables of a Word
func (w Word) Syllables() []string {
	sil := make([]string, 0)
	for _, s := range []string{w.Syllable4, w.Syllable3, w.Syllable2, w.Syllable1} {
		if s != "" {
			sil = append(sil, s)
		}
	}
	return sil
}

func (w Word) Sounds() []string {
	return soundsFromSyllables(w.Syllables())
}
