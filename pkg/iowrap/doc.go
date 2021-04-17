// Package iowrap intends to provide io wrappers that useful.

package iowrap

//go:generate go run github.com/golang/mock/mockgen -package iowrap -destination mock_test.go io Reader,Closer,ReaderAt,Seeker,Writer
