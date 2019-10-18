package types

// Capability is a uint64 which represents a capability storage service have.
type Capability uint64

// All capability that storage used.
const (
	CapabilitySegment Capability = 1 << iota
	CapabilityReach
)

// SegmentCapable returns whether this service support segment operations or not.
func (c Capability) SegmentCapable() bool {
	return c&CapabilitySegment == Capability(1)
}

// ReachAble returns whether this service support reach operations or not.
func (c Capability) ReachAble() bool {
	return c&CapabilityReach == Capability(1)
}
