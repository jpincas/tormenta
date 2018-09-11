package tormentadb_test

import (
	"testing"

	tormenta "github.com/jpincas/tormenta/tormentadb"
	"github.com/jpincas/tormenta/tormentadb/types"
)

func Test_Context(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	entity := types.TestType{}
	db.Save(&entity)

	sessionID := "session1234"

	db.First(&entity).SetContext("sessionid", sessionID).Run()
	if entity.TriggerString != sessionID {
		t.Errorf("Context was not set correctly.  Expecting: %s; Got: %s", sessionID, entity.TriggerString)
	}
}
