package types

// Filesystem captures size information for a container's filesystem and root filesystem
type Filesystem struct {
	SizeRw     int64 `json:"sizerw"`
	SizeRootFs int64 `json:"sizerootfs"`
}
