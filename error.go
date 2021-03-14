package bloomfilter

import (
	"fmt"
)

var (
	ErrCapacity  error = fmt.Errorf("bloomFilter: capacity must be greater than 0")
	ErrErrorRate error = fmt.Errorf("bloomFilter: error rate must be between 0 and 1")
)
