package main

import (
	"testing"
)

func BenchmarkConnect(b *testing.B) {
	     for i := 0; i < b.N; i++ {
	   	Connect()
	}
}

/*

go test -bench=.

Results:

429           2679830 ns/op
PASS
ok      milvus  3.521s

Function Connect() called = 429 times in 3.521 seconds.
Average Connection Time   = 2.68 milliseconds

429 / 3.521 = 121.9 connections per second

On Average, it takes 2,679,830 nanoseconds to connect to Milvus.

1 ms = 1,000,000 nanoseconds
2,679,830 nanoseconds = 2.67983 milliseconds

*/
