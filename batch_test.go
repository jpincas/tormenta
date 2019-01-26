package tormenta

import "testing"

func Test_NoBatches(t *testing.T) {
	testCases := []struct {
		noEntities, batchSize, expected int
	}{
		{0, 0, 0},
		{0, 10, 0},
		{1, 10000, 1},
		{10, 0, 0},
		{10, 10, 1},
		{10000, 10000, 1},
		{50000, 10000, 5},
		{55000, 10000, 6},
	}

	for _, testCase := range testCases {
		result := noBatches(testCase.noEntities, testCase.batchSize)
		if result != testCase.expected {
			t.Errorf("Testing batch sizer. Expecting %v, got %v", testCase.expected, result)
		}
	}
}

func Test_BatchStartEnd(t *testing.T) {
	// Remember - Golang slice endpoint is EXCLUSIVE

	testCases := []struct {
		testName                       string
		counter, batchSize, noEntities int
		expectedStart, expectedEnd     int
	}{
		{"all zero", 0, 0, 0, 0, 0},
		{"no entities, positive batch size", 0, 10, 0, 0, 0},
		{"no entities, positive batch size, counter", 1, 10, 0, 0, 0},
		{"1 entity", 0, 10000, 1, 0, 1},
		{"same batch size as entities", 0, 10, 10, 0, 10},
		{"batch size smaller than number of entities", 0, 9, 10, 0, 9},
		{"batch size larger than number of entities", 0, 20, 10, 0, 10},
		{"second batch, exact fit", 1, 10, 20, 10, 20},
		{"second batch, entities left over", 1, 10, 25, 10, 20},
		{"third batch, entities left over from last batch", 2, 10, 25, 20, 25},
	}

	for _, testCase := range testCases {
		a, b := batchStartAndEnd(testCase.counter, testCase.batchSize, testCase.noEntities)
		if a != testCase.expectedStart || b != testCase.expectedEnd {
			t.Errorf("Testing batch start/end calculator with %s. Expecting %v/%v, got %v/%v", testCase.testName, testCase.expectedStart, testCase.expectedEnd, a, b)
		}
	}
}
