// Package bitflags provides atomic flag registers and utilities for managing
package bitflags

import (
	"encoding"
	"fmt"
	"strings"
	"sync/atomic"
)

// Generic, atomic flag registers ---------------------------------------------

// FlagReg32 is an atomic register for bit flags with an underlying uint32.
// T must be a defined type whose underlying type is uint32.
type FlagReg32[T ~uint32] struct{ v atomic.Uint32 }

func (r *FlagReg32[T]) Load() T          { return T(r.v.Load()) }
func (r *FlagReg32[T]) Store(x T)        { r.v.Store(uint32(x)) }
func (r *FlagReg32[T]) Has(f T) bool     { return r.v.Load()&uint32(f) != 0 }
func (r *FlagReg32[T]) Any(mask T) bool  { return r.Has(mask) }
func (r *FlagReg32[T]) None(mask T) bool { return r.v.Load()&uint32(mask) == 0 }
func (r *FlagReg32[T]) Set(f T)          { r.cas(func(old uint32) uint32 { return old | uint32(f) }) }
func (r *FlagReg32[T]) Clear(f T)        { r.cas(func(old uint32) uint32 { return old &^ uint32(f) }) }
func (r *FlagReg32[T]) Toggle(f T)       { r.cas(func(old uint32) uint32 { return old ^ uint32(f) }) }
func (r *FlagReg32[T]) Mask(mask T) T    { return T(r.v.Load() & uint32(mask)) }
func (r *FlagReg32[T]) SetMask(mask, value T) {
	r.cas(func(old uint32) uint32 { return (old &^ uint32(mask)) | (uint32(value) & uint32(mask)) })
}
func (r *FlagReg32[T]) CompareAndSwap(old, new T) bool {
	return r.v.CompareAndSwap(uint32(old), uint32(new))
}
func (r *FlagReg32[T]) cas(f func(old uint32) uint32) {
	for {
		o := r.v.Load()
		n := f(o)
		if r.v.CompareAndSwap(o, n) {
			return
		}
	}
}

// Conditional helpers ---------------------------------------------------------

// SetIf sets bits in `set` only if *none* of bits in `mustBeClear` are present.
// Returns true when applied.
func (r *FlagReg32[T]) SetIf(mustBeClear, set T) bool {
	for {
		old := r.v.Load()
		if old&uint32(mustBeClear) != 0 {
			return false
		}
		newV := old | uint32(set)
		if r.v.CompareAndSwap(old, newV) {
			return true
		}
	}
}

// ClearIf clears bits in `clr` only if *all* bits in `mustBeSet` are present.
func (r *FlagReg32[T]) ClearIf(mustBeSet, clr T) bool {
	for {
		old := r.v.Load()
		if old&uint32(mustBeSet) != uint32(mustBeSet) {
			return false
		}
		newV := old &^ uint32(clr)
		if r.v.CompareAndSwap(old, newV) {
			return true
		}
	}
}

// FlagReg64 mirrors FlagReg32 for 64-bit sets.
type FlagReg64[T ~uint64] struct{ v atomic.Uint64 }

func (r *FlagReg64[T]) Load() T       { return T(r.v.Load()) }
func (r *FlagReg64[T]) Store(x T)     { r.v.Store(uint64(x)) }
func (r *FlagReg64[T]) Has(f T) bool  { return r.v.Load()&uint64(f) != 0 }
func (r *FlagReg64[T]) Set(f T)       { r.cas(func(o uint64) uint64 { return o | uint64(f) }) }
func (r *FlagReg64[T]) Clear(f T)     { r.cas(func(o uint64) uint64 { return o &^ uint64(f) }) }
func (r *FlagReg64[T]) Mask(mask T) T { return T(r.v.Load() & uint64(mask)) }
func (r *FlagReg64[T]) SetMask(mask, value T) {
	r.cas(func(o uint64) uint64 { return (o &^ uint64(mask)) | (uint64(value) & uint64(mask)) })
}
func (r *FlagReg64[T]) CompareAndSwap(old, new T) bool {
	return r.v.CompareAndSwap(uint64(old), uint64(new))
}
func (r *FlagReg64[T]) cas(f func(uint64) uint64) {
	for {
		o := r.v.Load()
		n := f(o)
		if r.v.CompareAndSwap(o, n) {
			return
		}
	}
}

// Pretty-print helpers (string encoding of flags) -----------------------------

type nameVal32 struct {
	name string
	val  uint32
}

type flagStringer32[T ~uint32] struct{ table []nameVal32 }

func NewStringer32[T ~uint32](pairs map[string]T) encoding.TextMarshaler {
	arr := make([]nameVal32, 0, len(pairs))
	for k, v := range pairs {
		arr = append(arr, nameVal32{name: k, val: uint32(v)})
	}
	return flagStringer32[T]{table: arr}
}

func (s flagStringer32[T]) MarshalText() ([]byte, error) {
	return []byte("<flags-stringer>"), nil // placeholder; use FlagString() per value
}

// FlagString renders a pipe-separated list of set flag names.
func FlagString[T ~uint32](val T, pairs map[string]T) string {
	if len(pairs) == 0 {
		return fmt.Sprintf("0x%X", uint32(val))
	}
	names := make([]string, 0, len(pairs))
	for name, bit := range pairs {
		if uint32(val)&uint32(bit) != 0 {
			names = append(names, name)
		}
	}
	if len(names) == 0 {
		return "<none>"
	}
	return strings.Join(names, "|")
}
