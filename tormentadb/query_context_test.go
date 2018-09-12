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


// Essentially the same test as above but on an indexed match query, this failed previously because an indexed
// search used the Public 'query.Get' function which did not take a context as a parameter and therefore simply
// passes the empty context to the PostGet trigger.
func Test_Context_Match(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	entity := types.TestType{}
	entity.IntField = 42
	db.Save(&entity)

	sessionID := "session1234"

	db.First(&entity).SetContext("sessionid", sessionID).Match("IntField", 42).Run()
	if entity.TriggerString != sessionID {
		t.Errorf("Context was not set correctly.  Expecting: %s; Got: %s", sessionID, entity.TriggerString)
	}
}


func Test_Context_Get(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	savedEntity := types.TestType{}
	db.Save(&savedEntity)

	entity := types.TestType{}
	entity.ID = savedEntity.ID

	sessionID := "session1234"
	ctx := make(map[string]interface{})
	ctx["sessionid"] =  sessionID

	db.GetWithContext(&entity, ctx)
	if entity.TriggerString != sessionID {
		t.Errorf("Context was not set correctly.  Expecting: %s; Got: %s", sessionID, entity.TriggerString)
	}
}

