package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
	// "time"
)

type TelemetryResult struct {
	address string
	balance string
	id      int
}

func main() {
	// Populate addresses
	addresses := populateAddress(500_000_00)
	fmt.Println("Finished populating addresses")

	numWorkers := 200
	batchSize := 1000
	doShit(addresses, numWorkers, batchSize)
}

func doShit(addresses []string, numWorkers int, batchSize int) {
	g := new(errgroup.Group)
	// Create a context to cancel the workers in case of an error
	g, ctx := errgroup.WithContext(context.Background())

	// Create buffered channels for tasks and results
	tasks := make(chan []string, numWorkers)
	results := make(chan []TelemetryResult, numWorkers)

	// Fan-out: Create worker goroutines
	for w := 1; w <= numWorkers; w++ {
		func(w int) {
			g.Go(func() error {
				return worker(ctx, w, tasks, results)
			})
		}(w)
	}

	// Create a goroutine to send tasks to workers
	go func() {
		orchestrator(ctx, tasks, addresses, batchSize)
		close(tasks)
	}()

	// Create a goroutine to wait for all workers to finish
	// check if there is an error and close the results channel
	go func() {
		if err := g.Wait(); err == nil {
			fmt.Println("All workers have finalized")
		} else {
			fmt.Println("Error received: ", err)
		}
		close(results)
	}()

	// Process results as they come in
	finalizedResults := processResults(results)
	if g.Wait() != nil {
		err := g.Wait()
		fmt.Println("Context is cancelled we are destroying everything")
		fmt.Println(err)
	} else {
		fmt.Println("Completed Finalized results: ", len(finalizedResults))
	}
}

func populateAddress(amount int) []string {
	addresses := make([]string, amount)
	for i := 0; i < amount; i++ {
		addresses[i] = fmt.Sprintf("Item %d", i)
	}
	return addresses
}

// worker performs the task on jobs received and sends results to the results channel.
func worker(ctx context.Context, id int, tasks <-chan []string, results chan<- []TelemetryResult) error {
	for {
		select {
		case task, ok := <-tasks:
			if !ok {
				// fmt.Printf("Worker %d stopping due to channel closed\n", id)
				return nil // Channel closed, stop the worker
			}

			processResults, err := performTask(task, id)
			if err != nil {
				// fmt.Println("Error received: ", err)
				return err
			}
			results <- processResults
		case <-ctx.Done():
			// Context cancelled, stop the worker
			// fmt.Printf("Worker %d stopping due to cancellation\n", id)
			return nil
		}
	}
}

func performTask(task []string, id int) ([]TelemetryResult, error) {
	results := make([]TelemetryResult, len(task))
	for i := range task {
		// if task[i] == "Item 900000" {
		// 	fmt.Println("Someone is here")
		// 	// select {
		// 	// case errs <- fmt.Errorf("Worker %d received error processing address: %s", id, task[i]):
		// 	// default:
		// 	// }
		// 	return nil, fmt.Errorf("Worker %d received error processing address: %s", id, task[i])
		// }
		// time.Sleep(10 * time.Millisecond)
		// fmt.Printf("Worker %d processing address: %s\n", id, task[i])
		time.Sleep(10 * time.Nanosecond)
		resp := fmt.Sprintf("Worker %d processing address: %s\n", id, task[i])
		// create a slice of slice of strings
		results[i] = TelemetryResult{address: resp, balance: "Balance", id: 3}
	}
	return results, nil
}

func orchestrator(ctx context.Context, tasks chan<- []string, addresses []string, batchSize int) {
	var currentBatch []string
	for i, v := range addresses {
		if ctx.Err() != nil {
			if ctx.Err() == context.Canceled {
				fmt.Println("Context is cancelled")
			} else if ctx.Err() == context.DeadlineExceeded {
				fmt.Println("Deadline has been exceeded")
			}
			// If the context is already cancelled, stop sending tasks
			break
		}

		currentBatch = append(currentBatch, v)
		// Check if the current batch is filled or it's the last element.
		if (i+1)%batchSize == 0 || i == len(addresses)-1 {
			tasks <- currentBatch
			currentBatch = nil // Reset current batch
		}
	}
}

func processResults(results <-chan []TelemetryResult) []TelemetryResult {
	finalizedResults := make([]TelemetryResult, 0)
	for batchResults := range results {
		for i := range batchResults {
			// length := len(finalizedResults)
			// if length%1000 == 0 {
			// 	// fmt.Println("Finalized results: ", len(finalizedResults))
			// }
			finalizedResults = append(finalizedResults, batchResults[i])
		}
	}
	return finalizedResults
}
