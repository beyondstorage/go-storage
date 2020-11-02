package types

import (
	"fmt"
)

// Pair will store option for storage service.
type Pair struct {
	Key   string
	Value interface{}
}

func (p Pair) String() string {
	return fmt.Sprintf("%s: %v", p.Key, p.Value)
}

type PairPolicyAction = uint8

const (
	PairPolicyActionError PairPolicyAction = iota
	PairPolicyActionIgnore
)
