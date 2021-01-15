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

func (l ListMode) IsDir() bool {
	return l&ListModeDir != 0
}
func (l ListMode) IsPrefix() bool {
	return l&ListModePrefix != 0
}

func (l ListMode) IsPart() bool {
	return l&ListModePart != 0
}

func (l ListMode) IsBlock() bool {
	return l&ListModeBlock != 0
}
