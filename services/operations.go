package services

// All public op are listed here for references.
const (
	// Init related op
	OpNewServicer = "new_servicer"
	OpNewStorager = "new_storager"

	// Service related op
	OpList   = "list"
	OpGet    = "get"
	OpCreate = "create"
	OpDelete = "delete"

	// Storage related op
	OpListPrefix = "list_prefix"
	OpListDir    = "list_dir"
	OpRead       = "read"
	OpWrite      = "write"
	OpStat       = "stat"
	OpMetadata   = "metadata"

	// Extended op
	OpCopy        = "copy"
	OpMove        = "move"
	OpReach       = "reach"
	OpStatistical = "statistical"

	// Segment related op
	OpListPrefixSegments = "list_prefix_segments"
	OpInitIndexSegment   = "init_index_segment"
	OpWriteIndexSegment  = "write_index_segment"
	OpCompleteSegment    = "complete_segment"
	OpAbortSegment       = "abort_segment"
)
