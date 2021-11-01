// Package iowrap intends to provide io wrappers that useful.

package iowrap

import _ "github.com/golang/mock/mockgen/model"

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -package iowrap -destination mock_test.go io Reader,Closer,ReaderAt,Seeker,Writer
