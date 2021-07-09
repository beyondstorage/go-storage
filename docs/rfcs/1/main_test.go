package main

import (
	"errors"
	"testing"
)

type TestStorager1 interface {
}

type Copier interface {
	Copy()
}

type test1 struct {
}

func (t *test1) Copy() {}

func BenchmarkCopierInterface(b *testing.B) {
	v := TestStorager1(&test1{})
	for i := 0; i < b.N; i++ {
		if x, ok := v.(Copier); ok {
			x.Copy()
		}
	}
}

type TestStorager2 interface {
	Copy()
	CopyAble() bool
}

type test2 struct{}

func (t *test2) Copy() {}

func (t *test2) CopyAble() bool {
	return true
}

func BenchmarkCopyableFuncCall(b *testing.B) {
	v := TestStorager2(&test2{})
	for i := 0; i < b.N; i++ {
		if v.CopyAble() {
			v.Copy()
		}
	}
}

type TestStorager3 interface {
	Capability() uint64
	Copy()
}

type test3 struct {
}

func (t *test3) Copy() {}

func (t *test3) Capability() uint64 {
	return 1
}

func BenchmarkCopyCapability(b *testing.B) {
	v := TestStorager3(&test3{})
	for i := 0; i < b.N; i++ {
		if v.Capability()&1 == 1 {
			v.Copy()
		}
	}
}

type TestStorager4 interface {
	Copy() error
	CopyPanic()
}

type test4 struct {
}

func (t *test4) Copy() error {
	return errors.New("test")
}

func (t *test4) CopyPanic() {
	panic("test")
}

func BenchmarkError(b *testing.B) {
	v := TestStorager4(&test4{})
	for i := 0; i < b.N; i++ {
		err := v.Copy()
		if err != nil {
			continue
		}
	}
}

func BenchmarkPanic(b *testing.B) {
	v := TestStorager4(&test4{})
	for i := 0; i < b.N; i++ {
		func() {
			defer func() {
				_ = recover()
			}()
			v.CopyPanic()
		}()
	}
}
