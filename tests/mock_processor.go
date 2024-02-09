package tests

import (
	"micro-batcher/pkg/batcher"
)

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
