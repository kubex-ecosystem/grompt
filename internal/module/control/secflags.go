package control

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

/* ========= FLAGS ========= */

type SecFlag uint32

const (
	SecNone         SecFlag = 0
	SecAuth         SecFlag = 1 << iota // autenticação/autorização
	SecSanitize                         // sanitize em params/headers
	SecSanitizeBody                     // sanitize no body
)

func (f SecFlag) Has(mask SecFlag) bool { return f&mask == mask }
func (f SecFlag) Any(mask SecFlag) bool { return f&mask != 0 }
func (f SecFlag) With(mask SecFlag) SecFlag {
	return f | mask
}
func (f SecFlag) Without(mask SecFlag) SecFlag {
	return f &^ mask
}

// ordem determinística para log/telemetria
var secOrder = []struct {
	name string
	flag SecFlag
}{
	{"auth", SecAuth},
	{"sanitize", SecSanitize},
	{"sanitize_body", SecSanitizeBody},
}

func (f SecFlag) String() string {
	if f == SecNone {
		return "none"
	}
	var parts []string
	for _, it := range secOrder {
		if f.Has(it.flag) {
			parts = append(parts, it.name)
		}
	}
	if len(parts) == 0 {
		return fmt.Sprintf("unknown(0x%X)", uint32(f))
	}
	return strings.Join(parts, "|")
}

/* ======== REGISTRADOR ATÔMICO ======== */

type FlagReg32A[T ~uint32] struct{ v atomic.Uint32 }

// Set  CAS OR (não usa Add; evita somas indevidas quando máscara tem múltiplos bits)
func (r *FlagReg32A[T]) Set(mask T) {
	for {
		old := r.v.Load()
		newV := old | uint32(mask)
		if r.v.CompareAndSwap(old, newV) {
			return
		}
	}
}

// Clear CAS AND NOT
func (r *FlagReg32A[T]) Clear(mask T) {
	for {
		old := r.v.Load()
		newV := old &^ uint32(mask)
		if r.v.CompareAndSwap(old, newV) {
			return
		}
	}
}

// Load current value
func (r *FlagReg32A[T]) Load() T { return T(r.v.Load()) }

// SetIf CAS OR se todos os bits de mustHave estiverem setados
func (r *FlagReg32A[T]) SetIf(mask, mustHave T) bool {
	for {
		old := r.v.Load()
		if (old & uint32(mustHave)) != uint32(mustHave) {
			return false
		}
		newV := old | uint32(mask)
		if r.v.CompareAndSwap(old, newV) {
			return true
		}
	}
}

// ======== MAP LEGADO -> FLAGS ========

// FromLegacyMap converte mapa legado (ex: de config JSON) para flags
func FromLegacyMap(m map[string]bool) SecFlag {
	if m == nil {
		return SecNone
	}
	var f SecFlag
	if m["auth"] {
		f |= SecAuth
	}
	if m["sanitize"] {
		f |= SecSanitize
	}
	if m["sanitize_body"] || m["validateAndSanitizeBody"] {
		f |= SecSanitizeBody
	}
	// compat antiga
	if m["validateAndSanitize"] {
		f |= SecSanitize
	}
	return f
}

/* ======== JOB STATES (bitmask) ======== */

type JobFlagC uint32

const (
	JobNone     JobFlagC = 0
	JobRunningC JobFlagC = 1 << iota
	JobRetryingC
	JobCompletedC
	JobFailedC
	JobTimedOutC
)

func (f JobFlagC) Has(mask JobFlagC) bool { return f&mask == mask }

type JobStateS struct{ reg FlagReg32[uint32] }

func (s *JobStateS) Load() JobFlagC { return JobFlagC(s.reg.Load()) }

func (s *JobStateS) Start() {
	s.reg.Clear(uint32(JobCompletedC | JobFailedC | JobTimedOutC))
	s.reg.Set(uint32(JobRunningC))
}

func (s *JobStateS) Retry() {
	// só marca retry se estava running
	_ = s.reg.SetIf(uint32(JobRetryingC), uint32(JobRunningC))
}

func (s *JobStateS) Complete() {
	// terminal limpa outros
	s.reg.Clear(uint32(JobRunningC | JobRetryingC | JobFailedC | JobTimedOutC))
	s.reg.Set(uint32(JobCompletedC))
}

func (s *JobStateS) Fail() {
	s.reg.Clear(uint32(JobRunningC | JobRetryingC | JobCompletedC | JobTimedOutC))
	s.reg.Set(uint32(JobFailedC))
}

func (s *JobStateS) Timeout() {
	s.reg.Clear(uint32(JobRunningC | JobRetryingC | JobCompletedC | JobFailedC))
	s.reg.Set(uint32(JobTimedOutC))
}

// String returns a string representation of the JobFlagC flags that are set.
// The output is a sorted, pipe-separated list of flag names (e.g., "completed|running").
// If no flags are set, it returns "none".
func (f JobFlagC) String() string {
	order := []struct {
		n string
		b JobFlagC
	}{
		{"running", JobRunningC},
		{"retrying", JobRetryingC},
		{"completed", JobCompletedC},
		{"failed", JobFailedC},
		{"timeout", JobTimedOutC},
	}
	var on []string
	for _, it := range order {
		if f.Has(it.b) {
			on = append(on, it.n)
		}
	}
	if len(on) == 0 {
		return "none"
	}
	sort.Strings(on)
	return strings.Join(on, "|")
}

func TestFlagReg32SetAndClear(t *testing.T) {
	var r FlagReg32[uint32]
	r.Set(uint32(SecAuth | SecSanitize))
	if got := SecFlag(r.Load()); !got.Has(SecAuth | SecSanitize) {
		t.Fatalf("expected flags set, got %v", got)
	}
	r.Clear(uint32(SecSanitize))
	if got := SecFlag(r.Load()); got.Has(SecSanitize) {
		t.Fatalf("sanitize should be cleared, got %v", got)
	}
}

func TestFlagReg32Concurrent(t *testing.T) {
	var r FlagReg32[uint32]
	wg := sync.WaitGroup{}
	N := 1000
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				r.Set(uint32(SecAuth))
			} else {
				r.Set(uint32(SecSanitize))
			}
		}(i)
	}
	wg.Wait()
	got := SecFlag(r.Load())
	if !got.Has(SecAuth) || !got.Has(SecSanitize) {
		t.Fatalf("expected both bits set, got %v", got)
	}
}

func TestJobStateTransitions(t *testing.T) {
	var s JobState
	s.Start()

	JobRunning := JobFlag(1 << 0)
	JobRetrying := JobFlag(1 << 1)
	JobCompleted := JobFlag(1 << 2)
	JobFailed := JobFlag(1 << 3)
	JobTimedOut := JobFlag(1 << 4)

	if st := s.Load(); !st.Has(JobRunning) {
		t.Fatalf("start → running, got %v", st)
	}
	s.Retry()
	if st := s.Load(); !st.Has(JobRetrying) {
		t.Fatalf("retry flag missing, got %v", st)
	}
	s.Complete()
	if st := s.Load(); !st.Has(JobCompleted) || st.Has(JobRunning|JobRetrying) {
		t.Fatalf("complete should be terminal only, got %v", st)
	}
	s.Start()
	s.Fail()
	if st := s.Load(); !st.Has(JobFailed) || st.Has(JobRunning) {
		t.Fatalf("failed should be terminal only, got %v", st)
	}
	s.Timeout()
	if st := s.Load(); !st.Has(JobTimedOut) || st.Has(JobRunning) {
		t.Fatalf("timeout should be terminal only, got %v", st)
	}
}
