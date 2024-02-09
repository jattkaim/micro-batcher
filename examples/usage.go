package main

import (
	"fmt"
	"github.com/jattkaim/micro-batcher/pkg/batcher"
	"time"
)

func main() {
	basicProcessor := &BasicProcessor{}
	microBatcher := batcher.NewBatcher(basicProcessor, 50, 10*time.Second)

	go func() {
		for result := range microBatcher.Results {
			if result.Error != nil {
				fmt.Printf("Job %s failed: %v\n", result.JobID, result.Error)
			} else {
				fmt.Printf("Job %s succeeded with result: %v\n", result.JobID, result.Result)
			}
		}
	}()

	// runs simulation of jobs being added to the batch
	for i := 0; i < 100; i++ {
		job := batcher.Job{ID: fmt.Sprintf("job-%d", i+1), Payload: fmt.Sprintf("payload-%d", i+1)}
		microBatcher.Add(job)
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(20 * time.Millisecond)

	microBatcher.Shutdown()
}
