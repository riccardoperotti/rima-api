package models

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sondsFromSilables(t *testing.T) {
	cases := []struct {
		silables []string
		sounds   []string
	}{
		{
			silables: []string{"per", "ga", "mi", "no"},
			sounds:   []string{"er", "a", "i", "o"},
		},
		{
			silables: []string{"ca", "ri", "ño", ""},
			sounds:   []string{"a", "i", "o"},
		},
		{
			silables: []string{"bar", "co", "", ""},
			sounds:   []string{"ar", "o"},
		},
		{
			silables: []string{"sol", "", "", ""},
			sounds:   []string{"ol"},
		},
		// TODO: Add a looooooooooot more tests. Should include diphthong, triphthongs and all that stuff
	}

	for i, test := range cases {
		test := test
		t.Run(fmt.Sprintf("%d: %s", i, strings.Join(test.silables[:], ",")), func(t *testing.T) {
			t.Parallel()

			r := sondsFromSilables(test.silables)

			assert.Equal(t, test.sounds, r, "Return from sondsFromSilables should be what we expect")
		})
	}

}

func Test_buildRimasSearchQuery(t *testing.T) {

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
				Name:         "península",
				SilableCount: 4,
				Type:         "E",
				Silable4:     "pe",
				Silable3:     "nín",
				Silable2:     "su",
				Silable1:     "la",
				EndsWith:     "a",
				Rank:         9,
			},
			query:    esdrujulaQ + sufix,
			bindVals: []interface{}{"%a", "E", "4", "%u", "%ín", "península"},
		},
		{
			word: Word{
				Name:         "tanque",
				SilableCount: 2,
				Type:         "G",
				Silable4:     "",
				Silable3:     "",
				Silable2:     "tan",
				Silable1:     "que",
				EndsWith:     "e",
				Rank:         2,
			},
			query:    graveQ + sufix,
			bindVals: []interface{}{"%e", "G", "2", "%an", "tanque"},
		},
		{
			word: Word{
				Name:         "serás",
				SilableCount: 2,
				Type:         "A",
				Silable4:     "",
				Silable3:     "",
				Silable2:     "se",
				Silable1:     "rás",
				EndsWith:     "s",
				Rank:         17,
			},
			query:    agudaQ + sufix,
			bindVals: []interface{}{"%ás", "A", "2", "serás"},
		},
	}

	for i, test := range cases {
		test := test
		t.Run(fmt.Sprintf("%d: %s", i, test.word.Name), func(t *testing.T) {
			t.Parallel()

			q, bv := buildRimasSearchQuery(test.word)

			assert.Equal(t, q, test.query, "query from buildRimasSearchQuery should be what we expect")
			assert.Equal(t, bv, test.bindVals, "bindValues from buildRimasSearchQuery should be what we expect")
		})
	}

}
