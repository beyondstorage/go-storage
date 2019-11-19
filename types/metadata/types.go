/*
Package metadata intend to provide all available metadata.
*/
package metadata

// Metadata will hold storager or object's metadata.
//
//go:generate ../../internal/bin/metadata
type Metadata map[string]interface{}
