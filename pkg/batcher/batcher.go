package batcher

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID      string
	Payload interface{}
}

type JobResult struct {
	JobID  string
	Result interface{}
	Error  error
}
type BatchProcessor interface {
	Process(jobs []Job) ([]JobResult, error)
	Start() error
	Stop() error
}

type Batcher struct {
	processor     BatchProcessor
	batchSize     int
	flushInterval time.Duration
	jobs          chan Job
	shutdown      chan struct{}
	Wg            sync.WaitGroup
	Results       chan JobResult
}

// NewBatcher creates and returns a new Batcher allowing user to specify batchSize & flushInterval
func NewBatcher(processor BatchProcessor, batchSize int, flushInterval time.Duration) *Batcher {
	batcher := &Batcher{
		processor:     processor,
		batchSize:     batchSize,
		flushInterval: flushInterval,
		jobs:          make(chan Job, batchSize),
		shutdown:      make(chan struct{}),
		Results:       make(chan JobResult, batchSize),
	}
	batcher.start()
	return batcher
}

// start begins the batch processing / flushing routine.
func (batcher *Batcher) start() {
	batcher.Wg.Add(1)
	go func() {
		defer batcher.Wg.Done()
		flushTicker := time.NewTicker(batcher.flushInterval) // ticker set to flushInterval
		defer flushTicker.Stop()

		var batch []Job

		for {
			select {
			case job := <-batcher.jobs:
				batch = append(batch, job)
				if len(batch) >= batcher.batchSize {
					batcher.runProcessor(batch)
					batch = nil
				}
			case <-flushTicker.C:
				if len(batch) > 0 {
					batcher.runProcessor(batch)
					batch = nil
				}
			case <-batcher.shutdown:
				if len(batch) > 0 {
					batcher.runProcessor(batch) // shutdown after processing
				}
				return
			}
		}
	}()
}

// runProcessor calls the external BatchProcessor to run a batch of jobs.
func (batcher *Batcher) runProcessor(jobs []Job) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic in batch processing: %v\n", r)
			// add an error result in each job if hard panic
			for _, job := range jobs {
				batcher.Results <- JobResult{JobID: job.ID, Result: nil, Error: fmt.Errorf("panic in batch processing")}
			}
		}
	}()
	results, err := batcher.processor.Process(jobs)
	// if there is an error while processing the batch-processor, attach error to all jobs
	if err != nil {
		for _, job := range jobs {
			batcher.Results <- JobResult{JobID: job.ID, Error: err, Result: job.Payload}
		}
		return
	}
	// TODO assumption here is that the results object contains `payload` and `error` this should be typed in future
	for _, result := range results {
		batcher.Results <- result
	}
}

// Add a job to the Batcher and returns a JobResult with a jobId, this is a placeholder result
func (batcher *Batcher) Add(job Job) JobResult {
	batcher.jobs <- job
	return JobResult{JobID: job.ID}
}

// Shutdown will trigger the Batcher to finish processing but only when all jobs are processed.
func (batcher *Batcher) Shutdown() {
	close(batcher.shutdown)
	batcher.Wg.Wait() // Wait for the processing goroutine to finish
}
