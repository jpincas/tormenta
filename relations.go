package tormenta

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/jpincas/gouuidv6"
)

const (
	ErrIncorrectIDField            = "%s does not exist or is not an ID field"
	ErrNoRecords                   = "at least 1 record is needed in order to load relations"
	ErrRelationMustBeStructPointer = "relation must be a pointer to a struct"
)

// TODO: clean this up when all relational stuff is done

func LoadRelations(db *DB, fieldName string, entities ...Record) error {
	// We need at least 1 entity to make this work
	if len(entities) == 0 {
		return errors.New(ErrNoRecords)
	}

	// Get all the related IDs for all the entities passed in
	// list arrived deduped so no need to worry about that here

	// Use a map to build a set of IDs to avoid duplication
	idMap := map[gouuidv6.UUID]bool{}

	for _, entity := range entities {
		idfieldName := fieldName + "ID"
		field := recordValue(entity).FieldByName(idfieldName)
		id, ok := field.Interface().(gouuidv6.UUID)
		if !ok {
			fmt.Errorf(ErrIncorrectIDField, idfieldName)
		}

		idMap[id] = true
	}

	// Map -> List
	var ids []gouuidv6.UUID
	for k := range idMap {
		ids = append(ids, k)
	}

	// Now we have the IDs of the related entities we need to get,
	// we just have to work out what type we are getting.
	// Use the first record as an exemplar -
	// check that its a pointer, and if so that it points to a struct
	fieldValue := fieldValue(entities[0], fieldName)
	// if fieldValue.Kind() != reflect.Ptr {
	// 	return errors.New(ErrRelationMustBeStructPointer)
	// }

	if reflect.ValueOf(fieldValue).Kind() != reflect.Struct {
		return errors.New(ErrRelationMustBeStructPointer)
	}

	// Set up a new slice of the type we are getting
	// and use the multiple Get by ID api to grab all the
	// relations
	results := newSlice(fieldValue, len(ids))
	if _, err := db.GetIDs(results, ids...); err != nil {
		return err
	}

	// At this point, results is *[]WhateverTheEntityIs
	// We'll iterate it and turn it into a map of 'read only' records
	// That's because we don't have pointers, so they
	// don't fulfil the full 'Record' interface.
	// It doesn't matter though - all we need is to be able to extract the ID
	recordMap := map[gouuidv6.UUID]ReadOnlyRecord{}
	s := reflect.ValueOf(results).Elem()
	for i := 0; i < s.Len(); i++ {
		record := s.Index(i).Interface().(ReadOnlyRecord)
		recordMap[record.GetID()] = record
	}

	// Now we need to iterate through the entities, setting the relation from the map
	for i := range entities {
		idfieldName := fieldName + "ID"
		field := recordValue(entities[i]).FieldByName(idfieldName)
		id, ok := field.Interface().(gouuidv6.UUID)
		if ok {
			// Get the record from the record map - if its nil
			// don't worry, the relation will just be nil
			record := recordMap[id]
			recordValue(entities[i]).FieldByName(fieldName).Set(reflect.ValueOf(record))
		}
	}

	return nil
}
