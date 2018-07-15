package tormentadb

const (
	batchSize = 10000
)

// noBatches works out how many save batches are necessary given N number of entities
func noBatches(n int, bs int) int {
	// If either the batch size or number of entities is 0,
	// then no batches are needed
	if bs == 0 || n == 0 {
		return 0
	}

	if (n % bs) > 0 {
		return (n / bs) + 1
	}
	return n / bs
}

func batchStartAndEnd(counter, bs, n int) (int, int) {
	// No entities -> 0/0 start/end
	if n == 0 {
		return 0, 0
	}

	start := counter * bs
	end := start + bs
	if end > n {
		end = n
	}

	return start, end
}
