package models

import (
	"fmt"
	"regexp"
)

// These are building blocks for regexes
const (
	wovel            = "[aáeéoóiíuú]"
	diphthong        = "ai|ái|au|áu|ei|éi|eu|éu|io|ió|ou|óu|ia|iá|ua|uá|ie|ié|ue|ué|oi|uo|uó|ui|iu|ay|áy|ey|éy|oy|eá"
	triphthong       = "iai|iei|uai|uei|uau|iau|uay|uey"
	consonantSimple  = "[bcdfghjklmnñpqrstvwyxz]"
	consonantComplex = "bl|br|ch|cl|cr|dr|fl|fr|gl|gr|kr|ll|pl|pr|rr|tr|tt|ss|qu|gü"
)

var anyVowel = fmt.Sprintf("(:?(:?%s)|(:?%s)|%s)", triphthong, diphthong, wovel)
var anyConsonant = fmt.Sprintf("(:?(:?%s)|%s)", consonantComplex, consonantSimple)

// Compiled regexes
var startsWithConsonants = regexp.MustCompile(`^` + anyConsonant + `+`)

// sondsFromSilables parses the sound of a silable, which is
// basically everything from the first vowel on
func sondsFromSilables(silables []string) []string {
	var sounds []string
	for _, sil := range silables {
		if sil == "" {
			continue
		}
		s := startsWithConsonants.ReplaceAllString(sil, "")
		sounds = append(sounds, s)
	}
	return sounds
}

// buildRimasSearchQuery builds the sql query and bind values array following
// rules of (Spanish) rhyming based on the Word's type
//
// TODO: there needs to be a strict/exact mode option where rhymes are not defined by sound,
// but are instead matched silable by silable.
func buildRimasSearchQuery(w Word) (string, []interface{}) {

	sounds := w.Sounds()

	// In all cases, the last sound, the type and the number of silables should always match:
	query := "select palabra, rank from lexico where silaba1 like ? and tipo = ? and silabas = ?"

	// bind values must be of type []any
	bindVals := []interface{}{
		fmtLike(sounds[len(sounds)-1]),
		w.Type,
		fmt.Sprintf("%d", w.SilableCount),
	}

	// if Type is G (GRAVE) - accent on the second to last silable,
	// the sound of the 2nd silable should also match
	if w.Type == "G" {
		query += " and silaba2 like ?"
		bindVals = append(bindVals, fmtLike(sounds[len(sounds)-2]))
	}

	// if Type is E (ESDRUJULA) - accent on the third to last silable,
	// the sound of the 2nd AND the 3rd silabas should also match
	if w.Type == "E" {
		query += " and silaba2 like ? and silaba3 like ?"
		bindVals = append(bindVals, fmtLike(sounds[len(sounds)-2]), fmtLike(sounds[len(sounds)-3]))
	}

	// and, of course, do not include the word we're trying to match!
	query += " and palabra != ?"
	bindVals = append(bindVals, w.Name)

	return query, bindVals
}

// fmtLike prepends wildcard '%' to string
func fmtLike(v string) string {
	return "%" + v
}
