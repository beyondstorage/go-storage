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
		input  ObjectMode
		args   ObjectMode
		expect ObjectMode
	}{
		{"simple case", ModeRead, ModeAppend, ModeRead | ModeAppend},
		{"complex case", ModeDir, ModeLink | ModeRead, ModeDir | ModeLink | ModeRead},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			v.input.Add(v.args)
			assert.Equal(t, v.expect, v.input)
		})
	}
}

func TestObjectMode_Del(t *testing.T) {
	cases := []struct {
		name   string
		input  ObjectMode
		args   ObjectMode
		expect ObjectMode
	}{
		{"simple case", ModeRead | ModeAppend, ModeAppend, ModeRead},
		{"complex case", ModeDir, ModeRead, ModeDir},
	}

	for _, v := range cases {

		t.Run(v.name, func(t *testing.T) {
			v.input.Del(v.args)
			assert.Equal(t, v.expect, v.input)
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
}
