package tormenta

// func Test_BasicWhere(t *testing.T) {
// 	db, _ := OpenTest("data/tests")
// 	defer db.Close()

// 	order1 := Order{}
// 	order2 := Order{
// 		Customer: 99,
// 	}
// 	db.Save(&order1, &order2)

// 	orders := []Order{}
// 	n, _ := db.Find(&orders).Where(
// 		Filter{
// 			"customer",
// 			func(indexContent string) bool {
// 				return indexContent == "99"
// 			},
// 		},
// 	).Run()

// 	if len(orders) != 2 || n != 2 {
// 		t.Errorf("Testing querying with 2 entity saved. Expecting 2 entities - got %v/%v", len(orders), n)
// 	}

// }
