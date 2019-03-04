type getFloat64Result struct {
	id     gouuidv6.UUID
	result float64
	found  bool
	err    error
}

func (db DB) getIDsWithContextFloat64AtPath(txn *badger.Txn, record Record, ctx map[string]interface{}, slowSumPath []string, ids ...gouuidv6.UUID) (float64, error) {
	var sum float64
	ch := make(chan getFloat64Result)
	defer close(ch)

	var wg sync.WaitGroup

	for _, id := range ids {
		wg.Add(1)

		go func(thisID gouuidv6.UUID) {
			f, found, err := db.getFloat64AtPath(txn, record, ctx, thisID, slowSumPath)
			ch <- getFloat64Result{
				id:     thisID,
				result: f,
				found:  found,
				err:    err,
			}
		}(id)
	}

	var errorsList []error
	go func() {
		for getResult := range ch {
			if getResult.err != nil {
				errorsList = append(errorsList, getResult.err)
			} else if getResult.found {
				sum = sum + getResult.result
			}

			// Only signal to the wait group that a record has been fetched
			// at this point rather than the anonymous func above, otherwise
			// you tend to lose the last result
			wg.Done()
		}
	}()

	// Once all the results are in, we need to
	// sort them according to the original order
	// But we'll bail now if there were any errors
	wg.Wait()

	if len(errorsList) > 0 {
		return sum, errorsList[0]
	}

	return sum, nil
}

func (db DB) getFloat64AtPath(txn *badger.Txn, entity Record, ctx map[string]interface{}, id gouuidv6.UUID, slowSumPath []string) (float64, bool, error) {
	var result float64

	item, err := txn.Get(newContentKey(KeyRoot(entity), id).bytes())
	// We are not treating 'not found' as an actual error,
	// instead we return 'false' and nil (unless there is an actual error)
	if err == badger.ErrKeyNotFound {
		return result, false, nil
	} else if err != nil {
		return result, false, err
	}

	if err := item.Value(func(val []byte) error {
		result, err = jsonparser.GetFloat(val, slowSumPath...)
		return err
	}); err != nil {
		return result, false, err
	}

	return result, true, nil
}