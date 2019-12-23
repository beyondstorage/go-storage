package credential

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewStatic(t *testing.T) {
	ak, sk := uuid.New().String(), uuid.New().String()
	s := NewStatic(ak, sk)
	assert.Equal(t, ak, s.accessKey)
	assert.Equal(t, sk, s.secretKey)
}

func TestStatic_Value(t *testing.T) {
	ak, sk := uuid.New().String(), uuid.New().String()
	v := Static{
		accessKey: ak,
		secretKey: sk,
	}.Value()
	assert.Equal(t, ak, v.AccessKey)
	assert.Equal(t, sk, v.SecretKey)
}
