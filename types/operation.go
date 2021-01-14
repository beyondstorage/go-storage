package types

// ListMode is the type for list, underlying type is int.
type ListMode uint8

const (
	// ListTypeDir means this list will use dir type.
	ListModeDir ListMode = 1 << iota
	// ListTypePrefix means this list will use prefix type.
	ListModePrefix
	ListModePart
	ListModeBlock
)
