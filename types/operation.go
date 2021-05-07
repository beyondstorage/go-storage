package types

import "strings"

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

// String implement Stringer for ListMode.
//
// An object with Dir,Part will print like "dir|part"
func (o ListMode) String() string {
	s := make([]string, 0)
	if o.IsDir() {
		s = append(s, "dir")
	}
	if o.IsPrefix() {
		s = append(s, "prefix")
	}
	if o.IsPart() {
		s = append(s, "part")
	}
	if o.IsBlock() {
		s = append(s, "block")
	}
	return strings.Join(s, "|")
}

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
