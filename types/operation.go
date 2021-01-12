package types

// ListType is the type for list, underlying type is int.
type ListType int

const (
	// ListTypeDir means this list will use dir type.
	ListTypeDir    ListType = 1
	// ListTypePrefix means this list will use prefix type.
	// NOTE: It's possible for prefix list type to return dirs
	ListTypePrefix ListType = 2
)
