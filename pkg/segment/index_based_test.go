package segment

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexBasedPart_String(t *testing.T) {
	p := IndexBasedPart{
		Size:  20,
		Index: 30,
	}
	assert.Equal(t, "IndexBasedPart {Index: 30, Size: 20}", p.String())
}

func TestNewIndexBasedSegment(t *testing.T) {
	s := NewIndexBasedSegment("test", "xxxx")
	assert.Equal(t, "test", s.Path)
	assert.Equal(t, "xxxx", s.ID)
	assert.NotNil(t, s.p)
}

func TestIndexBasedSegment_String(t *testing.T) {
	s := NewIndexBasedSegment("test", "xxxx")
	assert.Equal(t, "IndexBasedSegment {ID: xxxx, Path: test}", s.String())
}

func TestIndexBasedSegment_InsertPart(t *testing.T) {
	type fields struct {
		TotalSize int64
		ID        string
		Parts     map[int]*IndexBasedPart
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
		{"middle part", fields{3, "", map[int]*IndexBasedPart{
			0: {0, 1},
			2: {2, 1},
		}}, args{1, 1}, 1, false, nil},
		{"last part", fields{3, "", map[int]*IndexBasedPart{
			0: {0, 1},
			1: {1, 1},
		}}, args{2, 1}, 2, false, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &IndexBasedSegment{
				Segment: NewSegment("", tt.fields.ID),
				p:       tt.fields.Parts,
			}
			if s.p == nil {
				s.p = make(map[int]*IndexBasedPart)
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

func TestIndexBasedSegment_SortedParts(t *testing.T) {
	tests := []struct {
		name          string
		parts         map[int]*IndexBasedPart
		expectedParts []*IndexBasedPart
	}{
		{
			"normal case",
			map[int]*IndexBasedPart{
				0: {0, 5},
				5: {5, 5},
			},
			[]*IndexBasedPart{
				&IndexBasedPart{0, 5},
				&IndexBasedPart{5, 5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &IndexBasedSegment{
				p: tt.parts,
			}

			gotParts := s.Parts()
			assert.EqualValues(t, tt.expectedParts, gotParts)
		})
	}
}

func TestNewIndexBasedSegments(t *testing.T) {
	s := NewIndexBasedSegments()
	assert.NotNil(t, s)
}

func TestIndexBasedSegments_Len(t *testing.T) {
	s := IndexBasedSegments{s: map[string]*IndexBasedSegment{
		"x": NewIndexBasedSegment("def", "abc"),
	}}
	assert.Equal(t, 1, s.Len())
}

func TestIndexBasedSegments_Delete(t *testing.T) {
	s := IndexBasedSegments{s: map[string]*IndexBasedSegment{
		"x": NewIndexBasedSegment("def", "abc"),
	}}

	s.Delete("x")
	s.Delete("not exist")
	assert.Equal(t, 0, len(s.s))
}

func TestIndexBasedSegments_Insert(t *testing.T) {
	s := IndexBasedSegments{s: map[string]*IndexBasedSegment{
		"x": NewIndexBasedSegment("def", "abc"),
	}}

	s.Insert(&IndexBasedSegment{
		Segment: NewSegment("q", "y"),
	})
	assert.Equal(t, 2, len(s.s))
}

func TestIndexBasedSegments_Get(t *testing.T) {
	s := IndexBasedSegments{s: map[string]*IndexBasedSegment{
		"abc": NewIndexBasedSegment("def", "abc"),
	}}

	seg, err := s.Get("abc")
	assert.NoError(t, err)
	assert.Equal(t, "abc", seg.ID)
	assert.Equal(t, "def", seg.Path)

	seg, err = s.Get("not exist")
	assert.True(t, errors.Is(err, ErrSegmentNotFound))
	assert.Nil(t, seg)
}
