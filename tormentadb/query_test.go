package tormentadb_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jpincas/gouuidv6"
	"github.com/jpincas/tormenta/demo"
	tormenta "github.com/jpincas/tormenta/tormentadb"
)

// Basic Queries

func Test_BasicQuery(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	// 1 order
	order1 := demo.Order{}
	db.Save(&order1)

	var orders []demo.Order
	n, err := db.Find(&orders).Run()

	if err != nil {
		t.Error("Testing basic querying - got error")
	}

	if len(orders) != 1 || n != 1 {
		t.Errorf("Testing querying with 1 entity saved. Expecting 1 entity - got %v/%v", len(orders), n)
	}

	orders = []demo.Order{}
	c, err := db.Find(&orders).Count()
	if c != 1 {
		t.Errorf("Testing count 1 entity saved. Expecting 1 - got %v", c)
	}

	// 2 orders
	order2 := demo.Order{}
	db.Save(&order2)

	orders = []demo.Order{}
	if n, _ := db.Find(&orders).Run(); n != 2 {
		t.Errorf("Testing querying with 2 entity saved. Expecting 2 entities - got %v", n)
	}

	if c, _ := db.Find(&orders).Count(); c != 2 {
		t.Errorf("Testing count 2 entities saved. Expecting 2 - got %v", c)
	}
	if order1.ID == order2.ID {
		t.Errorf("Testing querying with 2 entities saved. 2 entities saved both have same ID")
	}
	if orders[0].ID == orders[1].ID {
		t.Errorf("Testing querying with 2 entities saved. 2 results returned. Both have same ID")
	}

	// Limit
	orders = []demo.Order{}
	if n, _ := db.Find(&orders).Limit(1).Run(); n != 1 {
		t.Errorf("Testing querying with 2 entities saved + limit. Wrong number of results received")
	}

	// Reverse - simple, only tests number received
	orders = []demo.Order{}
	if n, _ := db.Find(&orders).Reverse().Run(); n != 2 {
		t.Errorf("Testing querying with 2 entities saved + reverse. Expected %v, got %v", 2, n)
	}

	// Reverse + Limit - simple, only tests number received
	orders = []demo.Order{}
	if n, _ := db.Find(&orders).Reverse().Limit(1).Run(); n != 1 {
		t.Errorf("Testing querying with 2 entities saved + reverse + limit. Expected %v, got %v", 1, n)
	}

}

func Test_BasicQuery_First(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	order1 := demo.Order{}
	order2 := demo.Order{}
	db.Save(&order1, &order2)

	var order demo.Order
	n, err := db.First(&order).Run()

	if err != nil {
		t.Error("Testing first - got error")
	}

	if n != 1 {
		t.Errorf("Testing first. Expecting 1 entity - got %v", n)
	}

	if order.ID.IsNil() {
		t.Errorf("Testing first. Nil ID retrieved")
	}

	if order.ID != order1.ID {
		t.Errorf("Testing first. Order IDs are not equal - wrong order retrieved")
	}

	// Test nothing found (impossible range)
	n, _ = db.First(&order).From(time.Now()).To(time.Now()).Run()
	if n != 0 {
		t.Errorf("Testing first when nothing should be found.  Got n = %v", n)
	}
}

func Test_BasicQuery_DateRange(t *testing.T) {
	// Create a list of orders over a date range
	var orders []tormenta.Tormentable
	dates := []time.Time{
		// Now
		time.Now(),

		// Over the last week
		time.Now().Add(-1 * 24 * time.Hour),
		time.Now().Add(-2 * 24 * time.Hour),
		time.Now().Add(-3 * 24 * time.Hour),
		time.Now().Add(-4 * 24 * time.Hour),
		time.Now().Add(-5 * 24 * time.Hour),
		time.Now().Add(-6 * 24 * time.Hour),
		time.Now().Add(-7 * 24 * time.Hour),

		// Specific years
		time.Date(2009, time.January, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2010, time.January, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2011, time.January, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2012, time.January, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2013, time.January, 1, 1, 0, 0, 0, time.UTC),
	}

	for _, date := range dates {
		orders = append(orders, &demo.Order{
			Model: tormenta.Model{
				ID: gouuidv6.NewFromTime(date),
			},
		})
	}

	// Save the orders
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()
	db.Save(orders...)

	// Also another entity, to make sure there is no crosstalk
	db.Save(&demo.Product{
		Code:          "001",
		Name:          "Computer",
		Price:         999.99,
		StartingStock: 50,
		Description:   demo.DefaultDescription})

	// Quick check that all orders have saved correctly
	var results []demo.Order
	n, _ := db.Find(&results).Run()

	if len(results) != len(orders) || n != len(orders) {
		t.Errorf("Testing range query. Haven't even got to ranges yet. Just basic query expected %v - got %v/%v", len(orders), len(results), n)
		t.FailNow()
	}

	// Range test cases
	testCases := []struct {
		testName  string
		from, to  time.Time
		expected  int
		includeTo bool
		limit     int
		reverse   bool
		offset    int
	}{
		{"from right now - no orders expected, no 'to'", time.Now(), time.Time{}, 0, false, 0, false, 0},
		{"from beginning of time - all orders should be included, no 'to'", time.Time{}, time.Time{}, len(orders), false, 0, false, 0},
		{"from beginning of time - offset 1", time.Time{}, time.Time{}, len(orders) - 1, false, 0, false, 1},
		{"from beginning of time - offset 2", time.Time{}, time.Time{}, len(orders) - 2, false, 0, false, 2},
		{"from 2014, no 'to'", time.Date(2014, time.January, 1, 1, 0, 0, 0, time.UTC), time.Time{}, 8, false, 0, false, 0},
		{"from 1 hour ago, no 'to'", time.Now().Add(-1 * time.Hour), time.Time{}, 1, false, 0, false, 0},
		{"from beginning of time to now - expect all", time.Time{}, time.Now(), len(orders), true, 0, false, 0},
		{"from beginning of time to 2014 - expect 5", time.Time{}, time.Date(2014, time.January, 1, 1, 0, 0, 0, time.UTC), 5, true, 0, false, 0},
		{"from beginning of time to an hour ago - expect all but 1", time.Time{}, time.Now().Add(-1 * time.Hour), len(orders) - 1, true, 0, false, 0},
		{"from beginning of time - limit 1", time.Time{}, time.Time{}, 1, false, 1, false, 0},
		{"from beginning of time - limit 10", time.Time{}, time.Time{}, 10, false, 10, false, 0},
		{"from beginning of time - limit 10 - offset 2 (shouldnt affect number of results)", time.Time{}, time.Time{}, 10, false, 10, false, 2},
		{"from beginning of time - limit more than there are", time.Time{}, time.Time{}, len(orders), false, 0, false, 0},
		{"reversed - from beginning of time", time.Time{}, time.Time{}, 0, false, 0, true, 0},
		{"reverse - from now - no to", time.Now(), time.Time{}, len(orders), false, 0, true, 0},
		{"reverse - from now to 2014 - expect 8", time.Now(), time.Date(2014, time.January, 1, 1, 0, 0, 0, time.UTC), 8, true, 0, true, 0},
		{"reverse - from now to 2014 - limit 5 - expect 5", time.Now(), time.Date(2014, time.January, 1, 1, 0, 0, 0, time.UTC), 5, true, 5, true, 0},
	}

	for _, testCase := range testCases {
		rangequeryResults := []demo.Order{}
		query := db.Find(&rangequeryResults).From(testCase.from)

		if testCase.includeTo {
			query = query.To(testCase.to)
		}

		if testCase.limit > 0 {
			query = query.Limit(testCase.limit)
		}

		if testCase.reverse {
			query = query.Reverse()
		}

		if testCase.offset > 0 {
			query = query.Offset(testCase.offset)
		}

		n, _ := query.Run()
		c, _ := query.Count()

		// Count should always equal number of results
		if c != n {
			t.Errorf("Testing %s. Number of results does not equal count. Count: %v, Results: %v", testCase.testName, c, n)
		}

		// Test number of records retrieved
		if n != testCase.expected {
			t.Errorf("Testing %s (number orders retrieved). Expected %v - got %v", testCase.testName, testCase.expected, n)
		}

		// Test Count
		if c != testCase.expected {
			t.Errorf("Testing %s (count). Expected %v - got %v", testCase.testName, testCase.expected, c)
		}

	}

}

// Index Queries

func Test_IndexQuery_Match_String(t *testing.T) {
	customers := []string{"jon", "jonathan", "pablo"}
	var orders []tormenta.Tormentable

	for i := 0; i < 100; i++ {
		orders = append(orders, &demo.Order{
			Customer: customers[i%len(customers)],
		})
	}

	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()
	db.Save(orders...)

	testCases := []struct {
		testName      string
		match         interface{}
		reverse       bool
		expected      int
		expectedError error
	}{
		{"blank string", nil, false, 0, errors.New(tormenta.ErrNilInputMatchIndexQuery)},
		{"blank string", "", false, 0, nil},
		{"should not match any", "nocustomerwiththisname", false, 0, nil},
		{"matches 1 exactly with no interference", "pablo", false, 33, nil},
		{"matches 1 exactly and 1 prefix", "jon", false, 34, nil},
		{"matches 1 exactly and has same prefix as other", "jonathan", false, 33, nil},

		// Reversed - shouldn't make any difference to N
		{"blank string - reversed", nil, true, 0, errors.New(tormenta.ErrNilInputMatchIndexQuery)},
		{"blank string - reversed", "", true, 0, nil},
		{"should not match any - reversed", "nocustomerwiththisname", true, 0, nil},
		{"matches 1 exactly with no interference - reversed", "pablo", true, 33, nil},
		{"matches 1 exactly and 1 prefix - reversed", "jon", true, 34, nil},
		{"matches 1 exactly and has same prefix as other - reversed", "jonathan", true, 33, nil},
	}

	for _, testCase := range testCases {
		results := []demo.Order{}

		q := db.Find(&results).Where("customer", testCase.match)
		if testCase.reverse {
			q.Reverse()
		}

		n, err := q.Run()

		if testCase.expectedError != nil && err == nil {
			t.Errorf("Testing %s. Expected error [%v] but got none", testCase.testName, testCase.expectedError)
		}

		if testCase.expectedError == nil && err != nil {
			t.Errorf("Testing %s. Didn't expect error [%v]", testCase.testName, err)
		}

		if n != testCase.expected {
			t.Errorf("Testing %s.  Expecting %v, got %v", testCase.testName, testCase.expected, n)
		}
	}
}

func Test_IndexQuery_Range(t *testing.T) {
	// Set up 100 orders with increasing department, customer and shipping fee
	// and save
	var orders []tormenta.Tormentable

	for i := 0; i < 100; i++ {
		orders = append(orders, &demo.Order{
			Department:  i + 1,
			Customer:    fmt.Sprintf("customer-%v", string((i%26)+65)),
			ShippingFee: float64(i) + 0.99,
		})
	}

	// Randomise order before saving,
	// to ensure save order is not affecting retrieval
	// in some roundabout way
	tormenta.RandomiseTormentables(orders)

	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()
	db.Save(orders...)

	testCases := []struct {
		testName      string
		indexName     string
		start, end    interface{}
		expected      int
		expectedError error
	}{
		// Non existent index
		{"non existent index - no range", "notanindex", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery)},
		{"non existent index", "notanindex", 1, 2, 0, nil},

		// Int
		{"integer - no range", "department", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery)},
		{"integer - from 1", "department", 1, nil, 100, nil},
		{"integer - from 2", "department", 2, nil, 99, nil},
		{"integer - from 50", "department", 50, nil, 51, nil},
		{"integer - 1 to 2", "department", 1, 2, 2, nil},
		{"integer - 50 to 59", "department", 50, 59, 10, nil},
		{"integer - 1 to 100", "department", 1, 100, 100, nil},
		{"integer - to 50", "department", nil, 50, 50, nil},

		// String
		{"string - no range", "customer", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery)},
		{"string", "customer", "customer", nil, 100, nil},
		{"string - from A", "customer", "customer-A", nil, 100, nil},
		{"string - from B", "customer", "customer-B", nil, 96, nil},
		{"string - from Z", "customer", "customer-Z", nil, 3, nil},
		{"string - from A to Z", "customer", "customer-A", "customer-Z", 100, nil},
		{"string - to Z", "customer", nil, "customer-Z", 100, nil},

		// Float
		{"float - no range", "shippingfee", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery)},
		{"float", "shippingfee", 0, nil, 100, nil},
		{"float", "shippingfee", 0.99, nil, 100, nil},
		{"float - from 1.99", "shippingfee", 1.99, nil, 99, nil},
		{"float - from 50.99", "shippingfee", 50.99, nil, 50, nil},
		{"float - from 99.99", "shippingfee", 99.99, nil, 1, nil},
		{"float - to 20.99", "shippingfee", nil, 20.99, 21, nil},
	}

	for _, testCase := range testCases {
		rangequeryResults := []demo.Order{}
		n, err := db.
			Find(&rangequeryResults).
			Where(testCase.indexName, testCase.start, testCase.end).
			Run()

		if testCase.expectedError != nil && err == nil {
			t.Errorf("Testing %s. Expected error [%v] but got none", testCase.testName, testCase.expectedError)
		}

		if testCase.expectedError == nil && err != nil {
			t.Errorf("Testing %s. Didn't expect error [%v]", testCase.testName, err)
		}

		// Check for correct number of returned results
		if n != testCase.expected {
			t.Errorf("Testing %s (number orders retrieved). Expected %v - got %v", testCase.testName, testCase.expected, n)
		}

		// Check each member of the results for nil ID, customer and shipping fee
		for i, order := range rangequeryResults {
			if order.ID.IsNil() {
				t.Errorf("Testing %s.  Order no %v has nil ID", testCase.testName, i)
			}

			if order.Department == 0 {
				t.Errorf("Testing %s.  Order no %v has 0 department", testCase.testName, i)
			}

			if order.Customer == "" {
				t.Errorf("Testing %s.  Order no %v has blank customer", testCase.testName, i)
			}

			if order.ShippingFee == 0.0 {
				t.Errorf("Testing %s.  Order no %v has 0 shipping fee", testCase.testName, i)
			}
		}
	}

}

// Helper for making groups of depatments
func getDept(i int) int {
	if i <= 10 {
		return 1
	} else if i <= 20 {
		return 2
	} else {
		return 3
	}
}
func Test_IndexQuery_Range_MultipleIndexMembers(t *testing.T) {
	var orders []tormenta.Tormentable

	for i := 1; i <= 30; i++ {
		order := &demo.Order{
			Department: getDept(i),
		}

		orders = append(orders, order)
	}

	tormenta.RandomiseTormentables(orders)

	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()
	db.Save(orders...)

	testCases := []struct {
		testName   string
		start, end interface{}
		expected   int
	}{
		{"all departments", 1, 3, 30},
		{"departments 1, 2", 1, 2, 20},
		{"department 1", 1, 1, 10},
	}

	for _, testCase := range testCases {
		rangequeryResults := []demo.Order{}
		n, _ := db.
			Find(&rangequeryResults).
			Where("department", testCase.start, testCase.end).
			Run()

		// Check for correct number of returned results
		if n != testCase.expected {
			t.Errorf("Testing %s (number orders retrieved). Expected %v - got %v", testCase.testName, testCase.expected, n)
		}
	}
}

func Test_Aggregation(t *testing.T) {
	var products []tormenta.Tormentable

	for i := 1; i <= 30; i++ {
		product := &demo.Product{
			Price:         float64(i),
			StartingStock: i,
		}

		products = append(products, product)
	}

	tormenta.RandomiseTormentables(products)

	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()
	db.Save(products...)

	results := []demo.Product{}
	var intSum int32
	var floatSum float64
	expected := 465

	// Int32

	_, err := db.Find(&results).Where("startingstock", 1, 30).Sum(&intSum)
	if err != nil {
		t.Error("Testing int32 agreggation.  Got error")
	}

	expectedIntSum := int32(expected)
	if intSum != expectedIntSum {
		t.Errorf("Testing int32 agreggation. Expteced %v, got %v", expectedIntSum, intSum)
	}

	// Float64

	_, err = db.Find(&results).Where("price", 1.00, 30.00).Sum(&floatSum)
	if err != nil {
		t.Error("Testing float64 agreggation.  Got error")
	}

	expectedFloatSum := float64(expected)
	if floatSum != expectedFloatSum {
		t.Errorf("Testing float64 agreggation. Expteced %v, got %v", expectedFloatSum, floatSum)
	}
}

func Test_IndexQuery_DateRange(t *testing.T) {
	var orders []tormenta.Tormentable

	for i := 1; i <= 30; i++ {
		order := &demo.Order{
			Model: tormenta.Model{
				ID: gouuidv6.NewFromTime(time.Date(2009, time.November, i, 23, 0, 0, 0, time.UTC)),
			},
			Department: getDept(i),
		}

		orders = append(orders, order)
	}

	tormenta.RandomiseTormentables(orders)

	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()
	db.Save(orders...)

	testCases := []struct {
		testName        string
		indexRangeStart interface{}
		addFrom, addTo  bool
		from, to        time.Time
		expected        int
		indexRangeEnd   interface{}
	}{
		// Exact match tests (indexRangeEnd is nil)
		{"match department 1 - no date restriction", 1, false, false, time.Time{}, time.Time{}, 10, nil},
		{"match department 1 - from beginning of time", 1, true, false, time.Time{}, time.Now(), 10, nil},
		{"match department 1 - from beginning of time to now", 1, true, true, time.Time{}, time.Now(), 10, nil},
		{"match department 1 - from now (no to)", 1, true, false, time.Now(), time.Time{}, 0, nil},
		{"match department 1 - from 1st Nov (no to)", 1, true, false, time.Date(2009, time.November, 1, 23, 0, 0, 0, time.UTC), time.Time{}, 10, nil},
		{"match department 1 - from 5th Nov", 1, true, false, time.Date(2009, time.November, 5, 23, 0, 0, 0, time.UTC), time.Time{}, 6, nil},
		{"match department 1 - from 1st-5th Nov", 1, true, true, time.Date(2009, time.November, 1, 23, 0, 0, 0, time.UTC), time.Date(2009, time.November, 5, 23, 0, 0, 0, time.UTC), 5, nil},

		// Range match tests
		{"departments 1-3 - no date restriction", 1, false, false, time.Time{}, time.Time{}, 30, 3},
		{"departments 1-3 - from beginning of time", 1, true, false, time.Time{}, time.Time{}, 30, 3},
		{"departments 1-3 - from beginning of time to now", 1, true, true, time.Time{}, time.Now(), 30, 3},
		{"departments 1-3 - from now (no to)", 1, true, false, time.Now(), time.Time{}, 0, 3},
		{"departments 1-3 - from 5th Nov", 1, true, false, time.Date(2009, time.November, 5, 23, 0, 0, 0, time.UTC), time.Time{}, 26, 3},
		{"departments 1-3 - from 5th Nov - 15th Nov", 1, true, true, time.Date(2009, time.November, 5, 23, 0, 0, 0, time.UTC), time.Date(2009, time.November, 15, 23, 0, 0, 0, time.UTC), 11, 3},
		{"departments 1-2 - from 5th Nov - 15th Nov", 1, true, true, time.Date(2009, time.November, 5, 23, 0, 0, 0, time.UTC), time.Date(2009, time.November, 15, 23, 0, 0, 0, time.UTC), 11, 2},
		{"departments 1-2 - from 5th Nov - 9th Nov", 1, true, true, time.Date(2009, time.November, 5, 23, 0, 0, 0, time.UTC), time.Date(2009, time.November, 9, 23, 0, 0, 0, time.UTC), 5, 2},
	}

	for _, testCase := range testCases {
		rangequeryResults := []demo.Order{}

		// If only indexRangeStart is specified then its an exact match search
		query := db.Find(&rangequeryResults)
		if testCase.indexRangeEnd == nil {
			query = query.Where("department", testCase.indexRangeStart)
		} else {
			query = query.Where("department", testCase.indexRangeStart, testCase.indexRangeEnd)
		}

		if testCase.addFrom {
			query = query.From(testCase.from)
		}

		if testCase.addTo {
			query = query.To(testCase.to)
		}

		n, _ := query.Run()

		// Check for correct number of returned results
		if n != testCase.expected {
			t.Errorf("Testing %s (number orders retrieved). Expected %v - got %v", testCase.testName, testCase.expected, n)
		}

	}
}
