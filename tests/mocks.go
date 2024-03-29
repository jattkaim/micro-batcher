package tests

import "github.com/jattkaim/micro-batcher/pkg/batcher"

// MockBatchProcessor is used to test batcher functionality with a mocked processor dependency
type MockBatchProcessor struct {
	ProcessFunc func(jobs []batcher.Job) ([]batcher.JobResult, error)
}

func (m *MockBatchProcessor) Process(jobs []batcher.Job) ([]batcher.JobResult, error) {
	if m.ProcessFunc != nil {
		return m.ProcessFunc(jobs)
	}
	return nil, nil
}

func (m *MockBatchProcessor) Start() error {
	return nil
}

func (m *MockBatchProcessor) Stop() error {
	return nil
}
