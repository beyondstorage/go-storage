package segment

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSegment_GetPartIndex(t *testing.T) {
	type fields struct {
		TotalSize int64
		ID        string
		Parts     []*Part
	}
	type args struct {
		p *Part
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantCur int
		hasErr  bool
		wantErr error
	}{
		{"first part", fields{1, "", nil}, args{&Part{0, 1}}, 0, false, nil},
		{"middle part", fields{3, "", []*Part{
			{0, 1},
			{2, 1},
		}}, args{&Part{1, 1}}, 1, false, nil},
		{"last part", fields{3, "", []*Part{
			{0, 1},
			{1, 1},
		}}, args{&Part{2, 1}}, 2, false, nil},
		{"intersected part", fields{10, "", []*Part{
			{0, 5},
			{5, 5},
		}}, args{&Part{0, 2}}, 0, true, ErrPartIntersected},
	}
	for _, tt := range tests {
		s := &Segment{
			TotalSize: tt.fields.TotalSize,
			ID:        tt.fields.ID,
			Parts:     tt.fields.Parts,
		}
		gotCur, err := s.GetPartIndex(tt.args.p)
		if tt.hasErr {
			assert.Error(t, err, tt.name)
			assert.True(t, errors.Is(err, tt.wantErr), tt.name)
		} else {
			assert.NoError(t, err, tt.name)
		}
		assert.Equal(t, tt.wantCur, gotCur, tt.name)
	}
}

func TestSegment_InsertPart(t *testing.T) {
	type fields struct {
		TotalSize int64
		ID        string
		Parts     []*Part
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
		{"first part", fields{1, "", nil}, args{&Part{0, 1}}, false, nil},
		{"middle part", fields{3, "", []*Part{
			{0, 1},
			{2, 1},
		}}, args{&Part{1, 1}}, false, nil},
		{"last part", fields{3, "", []*Part{
			{0, 1},
			{1, 1},
		}}, args{&Part{2, 1}}, false, nil},
		{"intersected part", fields{10, "", []*Part{
			{0, 5},
			{5, 5},
		}}, args{&Part{0, 2}}, true, ErrPartIntersected},
	}
	for _, tt := range tests {
		s := &Segment{
			TotalSize: tt.fields.TotalSize,
			ID:        tt.fields.ID,
			Parts:     tt.fields.Parts,
		}

		err := s.InsertPart(tt.args.p)
		if tt.hasErr {
			assert.Error(t, err, tt.name)
			assert.True(t, errors.Is(err, tt.wantErr), tt.name)
		} else {
			assert.NoError(t, err, tt.name)
		}
	}
}

func TestSegment_ValidateParts(t *testing.T) {
	type fields struct {
		TotalSize int64
		ID        string
		Parts     []*Part
	}
	tests := []struct {
		name    string
		fields  fields
		hasErr  bool
		wantErr error
	}{
		{"single part", fields{1, "", []*Part{
			{0, 1},
		}}, false, nil},
		{"missing part at middle", fields{3, "", []*Part{
			{0, 1},
			{2, 1},
		}}, true, ErrSegmentNotFulfilled},
		{"size not correct", fields{3, "", []*Part{
			{0, 1},
			{1, 1},
		}}, true, ErrSegmentSizeNotMatch},
		{"two part", fields{10, "", []*Part{
			{0, 5},
			{5, 5},
		}}, false, nil},
	}
	for _, tt := range tests {
		s := &Segment{
			TotalSize: tt.fields.TotalSize,
			ID:        tt.fields.ID,
			Parts:     tt.fields.Parts,
		}
		err := s.ValidateParts()
		if tt.hasErr {
			assert.Error(t, err, tt.name)
			assert.True(t, errors.Is(err, tt.wantErr), tt.name)
		} else {
			assert.NoError(t, err, tt.name)
		}
	}
}
