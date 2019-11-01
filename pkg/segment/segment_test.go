package segment

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSegment_InsertPart(t *testing.T) {
	type fields struct {
		TotalSize int64
		ID        string
		Parts     map[int64]*Part
	}
	type args struct {
		p *Part
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		hasErr  bool
		wantErr error
	}{
		{"first part", fields{1, "", nil}, args{&Part{0, 1, 0}}, false, nil},
		{"middle part", fields{3, "", map[int64]*Part{
			0: {0, 1, 0},
			2: {2, 1, 1},
		}}, args{&Part{1, 1, 0}}, false, nil},
		{"last part", fields{3, "", map[int64]*Part{
			0: {0, 1, 0},
			1: {1, 1, 1},
		}}, args{&Part{2, 1, 0}}, false, nil},
		{"insert do not check intersected part", fields{10, "", map[int64]*Part{
			0: {0, 5, 0},
			5: {5, 5, 1},
		}}, args{&Part{0, 2, 0}}, false, nil},
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

			_, err := s.InsertPart(tt.args.p)
			if tt.hasErr {
				assert.Error(t, err, tt.name)
				assert.True(t, errors.Is(err, tt.wantErr), tt.name)
			} else {
				assert.NoError(t, err, tt.name)
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
			2: {2, 1, 1},
		}}, true, ErrSegmentNotFulfilled},
		{"two part", fields{"", map[int64]*Part{
			0: {0, 5, 0},
			5: {5, 5, 1},
		}}, false, nil},
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
