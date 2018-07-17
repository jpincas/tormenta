package tormentarest

import tormenta "github.com/jpincas/tormenta/tormentadb"

func entityRoot(entity tormenta.Tormentable) string {
	return string(tormenta.KeyRoot(entity))
}
