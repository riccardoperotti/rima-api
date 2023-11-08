package models

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"
)

type WordModel struct{}

// Word represents the a row in the `words` table
type Word struct {
	Word          string
	SyllableCount int
	Type          string
	Syllable4     string
	Syllable3     string
	Syllable2     string
	Syllable1     string
	EndsWith      string
	Rank          int
}

// GetWord fetches this word's complete record from the `words` table
func (wm WordModel) GetWord(db *gorm.DB, word string) (Word, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	var w Word
	db.First(&w, "word = ?", word)
	if err := db.WithContext(ctx).First(&w, "word = ?", word).Error; err != nil {
		log.Printf("Error when fetching word:  %s", err.Error())
		return Word{}, err
	}
	return w, nil
}


// Syllables returns the non-empty syllables of a Word
func (w Word) Syllables() []string {
	sylls := make([]string, 0)
	for _, s := range []string{w.Syllable4, w.Syllable3, w.Syllable2, w.Syllable1} {
		if s != "" {
			sylls = append(sylls, s)
		}
	}
	return sylls
}

func (w Word) Sounds() []string {
	return soundsFromSyllables(w.Syllables())
}
