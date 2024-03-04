package main

import (
	"fmt"
	"testing"
	"time"
)

var testCases = []struct {
	name       string
	numWorkers int
	batchSize  int
}{
	// {"100/200", 100, 200},
	// {"100/200", 200, 400},
	{"50/100", 5000, 1000},
}

const AmountOfAddresses = 500_000

func BenchmarkMain(b *testing.B) {
	addresses := populateAddress(AmountOfAddresses)
	fmt.Println("Finished populating addresses")

	b.ResetTimer()
	// Start at 1, for valid EIP155, see regexEIP155 variable.
	for _, tc := range testCases {
		b.Run(fmt.Sprintf("Case %s", tc.name), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				b.StartTimer()
				doShit(addresses, tc.numWorkers, tc.batchSize)
				b.StopTimer()
			}
		})
	}
}

func BenchmarkDumbIterator(b *testing.B) {
	addresses := populateAddress(AmountOfAddresses)
	fmt.Println("Finished populating addresses")

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.StartTimer()
		dumbIterator(addresses)
		b.StopTimer()
	}
}

func dumbIterator(addresses []string) {
	results := make([]TelemetryResult, 0)
	for i := range addresses {
		time.Sleep(10 * time.Nanosecond)
		// length := len(results)
		// if length%1000 == 0 {
		// 	// fmt.Println("Finalized results: ", len(results))
		// }
		resp := fmt.Sprintf("Worker processing address: %s\n", addresses[i])
		telResult := TelemetryResult{
			address: resp,
			balance: "balance",
			id:      3,
		}
		results = append(results, telResult)
	}
	fmt.Println("Finalized results: ", len(results))
}
