package common

import (
	"sync"
)

type StepPollState struct {
	inProgress int
	totalCount int
	lock       sync.Mutex
}

func NewStepPollState(totalCount int) *StepPollState {
	return &StepPollState{
		inProgress: 0,
		totalCount: totalCount,
	}
}

func (s *StepPollState) GetTotalCount() int {
	return s.totalCount
}

func (s *StepPollState) MaxAvailable() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.totalCount - s.inProgress
}

func (s *StepPollState) AcquireMaxAvailable() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	available := s.totalCount - s.inProgress
	s.inProgress += available
	return available
}

func (s *StepPollState) Release(count int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.inProgress -= count
	if s.inProgress < 0 {
		s.inProgress = 0
	}
}
