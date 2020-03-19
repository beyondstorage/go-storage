package segment

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPart_String(t *testing.T) {
	p := Part{
		Size:  20,
		Index: 30,
	}
	assert.Equal(t, "Part {Index: 30, Size: 20}", p.String())
}

func TestNewSegment(t *testing.T) {
	s := NewSegment("test", "xxxx")
	assert.Equal(t, "test", s.Path)
	assert.Equal(t, "xxxx", s.ID)
	assert.NotNil(t, s.p)
}

func TestSegment_String(t *testing.T) {
	s := NewSegment("test", "xxxx")
	assert.Equal(t, "Segment {ID: xxxx, Path: test}", s.String())
}

func TestSegment_InsertPart(t *testing.T) {
	type fields struct {
		TotalSize int64
		ID        string
		Parts     map[int]*Part
	}
	type args struct {
		idx  int
		size int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		index   int
		hasErr  bool
		wantErr error
	}{
		{"first part", fields{1, "", nil}, args{0, 1}, 0, false, nil},
		{"middle part", fields{3, "", map[int]*Part{
			0: {0, 1},
			2: {2, 1},
		}}, args{1, 1}, 1, false, nil},
		{"last part", fields{3, "", map[int]*Part{
			0: {0, 1},
			1: {1, 1},
		}}, args{2, 1}, 2, false, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Segment{
				ID: tt.fields.ID,
				p:  tt.fields.Parts,
			}
			if s.p == nil {
				s.p = make(map[int]*Part)
			}

			gotPart, err := s.InsertPart(tt.args.idx, tt.args.size)
			if tt.hasErr {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.index, gotPart.Index)
			}
		})
	}
}

func TestSegment_SortedParts(t *testing.T) {
	tests := []struct {
		name          string
		parts         map[int]*Part
		expectedParts []*Part
	}{
		{
			"normal case",
			map[int]*Part{
				0: {0, 5},
				5: {5, 5},
			},
			[]*Part{
				{0, 5},
				{5, 5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Segment{
				p: tt.parts,
			}

			gotParts := s.Parts()
			assert.EqualValues(t, tt.expectedParts, gotParts)
		})
	}
}

func TestNewSegments(t *testing.T) {
	s := NewSegments()
	assert.NotNil(t, s)
}

func TestSegments_Len(t *testing.T) {
	s := Segments{s: map[string]*Segment{
		"x": {
			ID:   "abc",
			Path: "def",
		},
	}}
	assert.Equal(t, 1, s.Len())
}

func TestSegments_Delete(t *testing.T) {
	s := Segments{s: map[string]*Segment{
		"x": {
			ID:   "abc",
			Path: "def",
		},
	}}

	s.Delete("x")
	s.Delete("not exist")
	assert.Equal(t, 0, len(s.s))
}

func TestSegments_Insert(t *testing.T) {
	s := Segments{s: map[string]*Segment{
		"x": {
			ID:   "abc",
			Path: "def",
		},
	}}

	s.Insert(&Segment{
		ID:   "y",
		Path: "q",
	})
	assert.Equal(t, 2, len(s.s))
}

func TestSegments_Get(t *testing.T) {
	s := Segments{s: map[string]*Segment{
		"abc": {
			ID:   "abc",
			Path: "def",
		},
	}}

	seg, err := s.Get("abc")
	assert.NoError(t, err)
	assert.Equal(t, "abc", seg.ID)
	assert.Equal(t, "def", seg.Path)

	seg, err = s.Get("not exist")
	assert.True(t, errors.Is(err, ErrSegmentNotFound))
	assert.Nil(t, seg)
}
