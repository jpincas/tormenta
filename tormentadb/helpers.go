package tormentadb

import "math/rand"

func RandomiseTormentables(slice []Tormentable) {
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func MemberString(valid []string, target string) bool {
	for _, validOption := range valid {
		if target == validOption {
			return true
		}
	}
	return false
}

var nonContentWords = []string{"on", "at", "the", "in", "a"}

func removeNonContentWords(strings []string) (results []string) {
	for _, s := range strings {
		if !MemberString(nonContentWords, s) {
			results = append(results, s)
		}
	}

	return
}
