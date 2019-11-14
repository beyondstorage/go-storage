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

func (t *test1) Copy() {
	return
}

func BenchmarkCopierInterface(b *testing.B) {
	var v TestStorager1
	v = &test1{}
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

type test2 struct {
}

func (t *test2) Copy() {
	return
}

func (t *test2) CopyAble() bool {
	return true
}

func BenchmarkCopyableFuncCall(b *testing.B) {
	var v TestStorager2
	v = &test2{}
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

func (t *test3) Copy() {
	return
}

func (t *test3) Capability() uint64 {
	return 1
}

func BenchmarkCopyCapability(b *testing.B) {
	var v TestStorager3
	v = &test3{}
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
	var v TestStorager4
	v = &test4{}
	for i := 0; i < b.N; i++ {
		err := v.Copy()
		if err != nil {
			continue
		}
	}
}

func BenchmarkPanic(b *testing.B) {
	var v TestStorager4
	v = &test4{}
	for i := 0; i < b.N; i++ {
		func() {
			defer func() {
				recover()
			}()
			v.CopyPanic()
		}()
	}
}
