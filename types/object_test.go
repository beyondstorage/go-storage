package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestObjectMode_String(t *testing.T) {
	cases := []struct {
		name   string
		input  ObjectMode
		expect string
	}{
		{"simple case", ModeDir, "dir"},
		{"complex case", ModeDir | ModeRead | ModeLink, "dir|read|link"},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			assert.Equal(t, v.expect, v.input.String())
		})
	}
}

func TestObjectMode_IsDir(t *testing.T) {
	cases := []struct {
		name   string
		input  ObjectMode
		expect bool
	}{
		{"simple case", ModeDir, true},
		{"complex case", ModeDir | ModeLink, true},
		{"negative case", ModeRead | ModeLink, false},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			assert.Equal(t, v.expect, v.input.IsDir())
		})
	}
}

func TestObjectMode_Add(t *testing.T) {
	cases := []struct {
		name   string
		base   ObjectMode
		input  ObjectMode
		expect ObjectMode
	}{
		{"add single new mode", ModeRead, ModeAppend, ModeRead | ModeAppend},
		{"add multiple new modes", ModeDir, ModeLink | ModeRead, ModeDir | ModeLink | ModeRead},
		{"add its own", ModeRead, ModeRead, ModeRead},
		{"add exist mode", ModeRead, ModeAppend | ModeRead, ModeRead | ModeAppend},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			v.base.Add(v.input)
			assert.Equal(t, v.expect, v.base)
		})
	}
}

func TestObjectMode_Del(t *testing.T) {
	cases := []struct {
		name   string
		base   ObjectMode
		input  ObjectMode
		expect ObjectMode
	}{
		{"delete exist mode", ModeRead | ModeAppend, ModeAppend, ModeRead},
		{"delete absent mode", ModeDir, ModeRead, ModeDir},
		{"delete its own", ModeRead, ModeRead, 0},
		{"delete exist and absent mode", ModeRead | ModeDir, ModeRead | ModeAppend, ModeDir},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			v.base.Del(v.input)
			assert.Equal(t, v.expect, v.base)
		})
	}
}

func ExampleObjectMode_Add() {
	var o ObjectMode
	o.Add(ModeDir)
	o.Add(ModeRead | ModeAppend)
}

func ExampleObjectMode_Del() {
	o := ModeRead | ModeAppend
	o.Del(ModeAppend)
	o.Del(ModeRead | ModeAppend)
}
