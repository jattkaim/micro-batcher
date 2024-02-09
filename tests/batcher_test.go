package tests

import (
	"github.com/jattkaim/micro-batcher/pkg/batcher"
	"testing"
	"time"
)

func TestProcess(t *testing.T) {
	mp := &MockBatchProcessor{ProcessFunc: func(jobs []batcher.Job) ([]batcher.JobResult, error) {
		if len(jobs) != 2 {
			t.Errorf("Expected batch size of 2, got %d", len(jobs))
		}
		return make([]batcher.JobResult, len(jobs)), nil
	},
	}

	b := batcher.NewBatcher(mp, 2, 1*time.Second)

	b.Add(batcher.Job{ID: "a"})
	b.Add(batcher.Job{ID: "b"})

	b.Shutdown()
}

func TestProcessFailure(t *testing.T) {
	mp := &MockBatchProcessor{ProcessFunc: func(jobs []batcher.Job) ([]batcher.JobResult, error) {
		panic("hard failure in processing")
	},
	}

	b := batcher.NewBatcher(mp, 2, 1*time.Second)

	b.Add(batcher.Job{ID: "a"})
	b.Add(batcher.Job{ID: "b"})

	b.Shutdown()
}

func TestFlushInterval(t *testing.T) {
	processed := false
	mp := &MockBatchProcessor{ProcessFunc: func(jobs []batcher.Job) ([]batcher.JobResult, error) {
		processed = true
		return make([]batcher.JobResult, len(jobs)), nil
	},
	}
	b := batcher.NewBatcher(mp, 1, 100*time.Millisecond)
	b.Add(batcher.Job{ID: "a"})
	time.Sleep(110 * time.Millisecond)

	if !processed {
		t.Error("Expected jobs to be processed on flush interval")
	}

	b.Shutdown()
}

func TestShutdown(t *testing.T) {
	processed := false
	mp := &MockBatchProcessor{ProcessFunc: func(jobs []batcher.Job) ([]batcher.JobResult, error) {
		processed = true
		return make([]batcher.JobResult, len(jobs)), nil
	},
	}
	b := batcher.NewBatcher(mp, 10, 1*time.Second)
	b.Add(batcher.Job{ID: "a"})

	// following go routine starts the shutdown process after 50 milliseconds
	go func() {
		time.Sleep(50 * time.Millisecond)
		b.Shutdown()
	}()

	// TODO The wait group is exposed here to ensure the shutdown has completed before checking flag
	// TODO unsure if this is correct.
	b.Wg.Wait()

	if !processed {
		t.Error("Expected jobs to be processed before shutdown")
	}
}
