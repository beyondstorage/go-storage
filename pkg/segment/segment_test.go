package segment

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSegment(t *testing.T) {
	s := NewSegment("test", "xxxx", 10)
	assert.Equal(t, "test", s.Path)
	assert.Equal(t, "xxxx", s.ID)
	assert.Equal(t, int64(10), s.PartSize)
	assert.NotNil(t, s.Parts)
}

func TestSegment_InsertPart(t *testing.T) {
	type fields struct {
		TotalSize int64
		ID        string
		Parts     map[int64]*Part
	}
	type args struct {
		offset int64
		size   int64
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
		{"middle part", fields{3, "", map[int64]*Part{
			0: {0, 1, 0},
			2: {2, 1, 2},
		}}, args{1, 1}, 1, false, nil},
		{"last part", fields{3, "", map[int64]*Part{
			0: {0, 1, 0},
			1: {1, 1, 1},
		}}, args{2, 1}, 2, false, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Segment{
				ID:       tt.fields.ID,
				Parts:    tt.fields.Parts,
				PartSize: 1,
			}
			if s.Parts == nil {
				s.Parts = make(map[int64]*Part)
			}

			gotPart, err := s.InsertPart(tt.args.offset, tt.args.size)
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

func TestSegment_ValidateParts(t *testing.T) {
	type fields struct {
		ID    string
		Parts map[int64]*Part
	}
	tests := []struct {
		name    string
		fields  fields
		hasErr  bool
		wantErr error
	}{
		{"single part", fields{"", map[int64]*Part{
			0: {0, 1, 0},
		}}, false, nil},
		{"missing part at middle", fields{"", map[int64]*Part{
			0: {0, 1, 0},
			2: {2, 1, 2},
		}}, true, ErrSegmentNotFulfilled},
		{"two part", fields{"", map[int64]*Part{
			0: {0, 5, 0},
			5: {5, 5, 1},
		}}, false, nil},
		{"empty parts", fields{"", map[int64]*Part{}}, true, ErrSegmentPartsEmpty},
		{"first part is not 0", fields{"", map[int64]*Part{
			1: {1, 5, 0},
		}}, true, ErrSegmentNotFulfilled},
		{"intersected part", fields{"", map[int64]*Part{
			0: {0, 5, 0},
			2: {2, 5, 1},
		}}, true, ErrPartIntersected},
		{"not fulfilled part", fields{"", map[int64]*Part{
			0:  {0, 5, 0},
			10: {10, 5, 2},
		}}, true, ErrSegmentNotFulfilled},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Segment{
				ID:    tt.fields.ID,
				Parts: tt.fields.Parts,
			}
			if s.Parts == nil {
				s.Parts = make(map[int64]*Part)
			}

			err := s.ValidateParts()
			if tt.hasErr {
				assert.Error(t, err, tt.name)
				assert.True(t, errors.Is(err, tt.wantErr), tt.name)
			} else {
				assert.NoError(t, err, tt.name)
			}
		})
	}
}

func TestSegment_SortedParts(t *testing.T) {
	tests := []struct {
		name          string
		parts         map[int64]*Part
		expectedParts []*Part
	}{
		{
			"normal case",
			map[int64]*Part{
				0: {0, 5, 0},
				5: {5, 5, 1},
			},
			[]*Part{
				{0, 5, 0},
				{5, 5, 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Segment{
				Parts: tt.parts,
			}

			gotParts := s.SortedParts()
			assert.EqualValues(t, tt.expectedParts, gotParts)
		})
	}
}
