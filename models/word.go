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
	Name         string
	SilableCount int
	Type         string
	Silable4     string
	Silable3     string
	Silable2     string
	Silable1     string
	EndsWith     string
	Rank         int
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
		&dbW.SilableCount,
		&dbW.Type,
		&dbW.Silable4,
		&dbW.Silable3,
		&dbW.Silable2,
		&dbW.Silable1,
		&dbW.EndsWith,
		&dbW.Rank,
	)
	if err != nil {
		log.Printf("Error when scanning results:  %s", err.Error())
		return Word{}, err
	}

	return dbW, nil
}

// Silables returns the non-empty silables of a Word
func (w Word) Silables() []string {
	sil := make([]string, 0)
	for _, s := range []string{w.Silable4, w.Silable3, w.Silable2, w.Silable1} {
		if s != "" {
			sil = append(sil, s)
		}
	}
	return sil
}

func (w Word) Sounds() []string {
	return sondsFromSilables(w.Silables())
}
