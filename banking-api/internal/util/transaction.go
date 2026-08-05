package util

import (
	"fmt"
	"sync/atomic"
	"time"
)

var counter uint64

func NewTransactionID() string {

	id := atomic.AddUint64(&counter, 1)

	return fmt.Sprintf(
		"TX-%s-%06d",
		time.Now().Format("20060102"),
		id,
	)
}
