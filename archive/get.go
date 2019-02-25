package archive

import (
	"reflect"

	"github.com/jpincas/gouuidv6"
)

// For benchmarking / comparison with parallel get
func (db DB) GetIDsSerial(target interface{}, ids ...gouuidv6.UUID) (int, error) {
	records := newResultsArray(target)

	var counter int
	for _, id := range ids {
		// It's inefficient creating a new entity target for the result
		// on every loop, but we can't just create a single one
		// and reuse it, because there would be risk of data from 'previous'
		// entities 'infecting' later ones if a certain field wasn't present
		// in that later entity, but was in the previous one.
		// Unlikely if the all JSON is saved with the schema, but I don't
		// think we can risk it
		record := newRecordFromSlice(target)

		// For an error, we'll bail, if we simply can't find the record, we'll continue
		if found, err := db.get(record, noCTX, id); err != nil {
			return counter, err
		} else if found {
			records = reflect.Append(records, recordValue(record))
			counter++
		}
	}

	return counter, nil
}
