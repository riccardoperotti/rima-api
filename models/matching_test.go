package models

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_soundsFromSyllables(t *testing.T) {
	cases := []struct {
		syllables []string
		sounds    []string
	}{
		{
			syllables: []string{"per", "ga", "mi", "no"},
			sounds:    []string{"er", "a", "i", "o"},
		},
		{
			syllables: []string{"ca", "ri", "ño", ""},
			sounds:    []string{"a", "i", "o"},
		},
		{
			syllables: []string{"bar", "co", "", ""},
			sounds:    []string{"ar", "o"},
		},
		{
			syllables: []string{"sol", "", "", ""},
			sounds:    []string{"ol"},
		},
		// TODO: Add a looooooooooot more tests. Should include diphthong, triphthongs and all that stuff
	}

	for i, test := range cases {
		test := test
		t.Run(fmt.Sprintf("%d: %s", i, strings.Join(test.syllables[:], ",")), func(t *testing.T) {
			t.Parallel()

			r := soundsFromSyllables(test.syllables)

			assert.Equal(t, test.sounds, r, "Return from soundsFromSyllables should be what we expect")
		})
	}

}

func Test_buildRhymesSearchQuery(t *testing.T) {

	agudaQ := "select palabra, rank from lexico where silaba1 like ? and tipo = ? and silabas = ?"
	graveQ := agudaQ + " and silaba2 like ?"
	esdrujulaQ := graveQ + " and silaba3 like ?"
	sufix := " and palabra != ?"

	cases := []struct {
		word     Word
		query    string
		bindVals []interface{}
	}{
		{
			word: Word{
				Name:          "península",
				SyllableCount: 4,
				Type:          "E",
				Syllable4:     "pe",
				Syllable3:     "nín",
				Syllable2:     "su",
				Syllable1:     "la",
				EndsWith:      "a",
				Rank:          9,
			},
			query:    esdrujulaQ + sufix,
			bindVals: []interface{}{"%a", "E", "4", "%u", "%ín", "península"},
		},
		{
			word: Word{
				Name:          "tanque",
				SyllableCount: 2,
				Type:          "G",
				Syllable4:     "",
				Syllable3:     "",
				Syllable2:     "tan",
				Syllable1:     "que",
				EndsWith:      "e",
				Rank:          2,
			},
			query:    graveQ + sufix,
			bindVals: []interface{}{"%e", "G", "2", "%an", "tanque"},
		},
		{
			word: Word{
				Name:          "serás",
				SyllableCount: 2,
				Type:          "A",
				Syllable4:     "",
				Syllable3:     "",
				Syllable2:     "se",
				Syllable1:     "rás",
				EndsWith:      "s",
				Rank:          17,
			},
			query:    agudaQ + sufix,
			bindVals: []interface{}{"%ás", "A", "2", "serás"},
		},
	}

	for i, test := range cases {
		test := test
		t.Run(fmt.Sprintf("%d: %s", i, test.word.Name), func(t *testing.T) {
			t.Parallel()

			q, bv := buildRhymesSearchQuery(test.word)

			assert.Equal(t, q, test.query, "query from buildRhymesSearchQuery should be what we expect")
			assert.Equal(t, bv, test.bindVals, "bindValues from buildRhymesSearchQuery should be what we expect")
		})
	}

}
