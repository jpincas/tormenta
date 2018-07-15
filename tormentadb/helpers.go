package tormentadb

import "math/rand"

func RandomiseTormentables(slice []Tormentable) {
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}
