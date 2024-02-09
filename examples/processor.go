package main

import (
	"fmt"
	"micro-batcher/pkg/batcher"
)

// BasicProcessor is a simple implementation of a processor interface that our batcher would depend on.
type BasicProcessor struct{}

// Process just logs the jobs it processes and marks them as processed.
// Note: This BasicProcessor was written with ChatGPT :p
func (bp *BasicProcessor) Process(jobs []batcher.Job) ([]batcher.JobResult, error) {
	results := make([]batcher.JobResult, 0, len(jobs))
	for _, job := range jobs {
		// generate a result
		result := batcher.JobResult{
			JobID:  job.ID,
			Result: "success",
			Error:  nil,
		}
		results = append(results, result)
	}
	return results, nil
}

func (bp *BasicProcessor) Start() error {
	// Optionally, include start logic here
	fmt.Println("BasicProcessor started")
	return nil
}

func (bp *BasicProcessor) Stop() error {
	// Optionally, include stop logic here
	fmt.Println("BasicProcessor stopped")
	return nil
}
