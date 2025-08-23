// Package control provides abstractions for managing job states and security flags.
package control

import (
	"errors"
	"sync/atomic"
)

// Example: Job state flags (combináveis) -------------------------------------

type JobFlag uint32

const (
	JobPendingA JobFlag = 1 << iota
	JobRunningA
	JobCancelRequestedA
	JobRetryingA
	JobCompletedA
	JobFailedA
	JobTimedOutA
)

const (
	terminalMask JobFlag = JobCompletedA | JobFailedA | JobTimedOutA
)

var (
	ErrTerminal = errors.New("job is in a terminal state")
)

func (j JobFlag) Has(flag JobFlag) bool {
	return j&flag != 0
}

type FlagReg32[T ~uint32] struct{ v atomic.Uint32 }

func (r *FlagReg32[T]) Load() T { return T(r.v.Load()) }
func (r *FlagReg32[T]) Store(val T) {
	r.v.Store(uint32(val))
}
func (r *FlagReg32[T]) Set(mask T) { r.v.Add(uint32(mask)) }
func (r *FlagReg32[T]) Clear(mask T) {
	for {
		old := r.v.Load()
		if r.v.CompareAndSwap(old, old&^uint32(mask)) {
			return
		}
	}
}
func (r *FlagReg32[T]) SetIf(clearMask, setMask T) bool {
	for {
		old := r.v.Load()
		if old&uint32(clearMask) != 0 {
			return false
		}
		newV := (old &^ uint32(clearMask)) | uint32(setMask)
		if r.v.CompareAndSwap(old, newV) {
			return true
		}
	}
}
func (r *FlagReg32[T]) Any(mask T) bool { return r.v.Load()&uint32(mask) != 0 }
func (r *FlagReg32[T]) All(mask T) bool { return r.v.Load()&uint32(mask) == uint32(mask) }

type JobState struct{ r FlagReg32[JobFlag] }

func (s *JobState) Load() JobFlag { return s.r.Load() }

// Start only from Pending; sets Running.
func (s *JobState) Start() error {
	ok := s.r.SetIf(terminalMask|JobRunningA|JobCompletedA|JobFailedA|JobTimedOutA, JobRunningA)
	if !ok {
		return ErrTerminal
	}
	return nil
}

// RequestCancel sets the CancelRequested flag.
func (s *JobState) RequestCancel() { s.r.Set(JobCancelRequestedA) }

// Retry sets the job state to Retrying if not in a terminal state.
func (s *JobState) Retry() error {
	// can retry if not terminal; set Retrying and clear Running
	for {
		old := s.r.Load()
		if old&terminalMask != 0 {
			return ErrTerminal
		}
		newV := (old | JobRetryingA) &^ JobRunningA
		if s.r.SetIf(old, newV) {
			return nil
		}
	}
}

// Complete sets the job state to Completed.
func (s *JobState) Complete() error {
	for {
		old := s.r.Load()
		if old&terminalMask != 0 {
			return ErrTerminal
		}
		newV := (old | JobCompletedA) &^ (JobRunningA | JobRetryingA | JobCancelRequestedA)
		if s.r.SetIf(old, newV) {
			return nil
		}
	}
}

// Fail sets the job state to Failed.
func (s *JobState) Fail() error {
	for {
		old := s.r.Load()
		if old&terminalMask != 0 {
			return ErrTerminal
		}
		newV := (old | JobFailedA) &^ (JobRunningA | JobRetryingA)
		if s.r.SetIf(old, newV) {
			return nil
		}
	}
}

// Timeout sets the job state to TimedOut.
func (s *JobState) Timeout() error {
	for {
		old := s.r.Load()
		if old&terminalMask != 0 {
			return ErrTerminal
		}
		newV := (old | JobTimedOutA) &^ (JobRunningA | JobRetryingA)
		if s.r.SetIf(old, newV) {
			return nil
		}
	}
}

// IsTerminal returns true if the job is in a terminal state.
func (s *JobState) IsTerminal() bool { return s.r.Any(terminalMask) }
