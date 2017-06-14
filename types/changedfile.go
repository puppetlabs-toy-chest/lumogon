package types

// ChangedFile represents a file in a container that has been changed
// from that originally provided in the containers image, where Path
// is the file which has changed and Kind is the type of change.
// This is based on Dockers container.ContainerChangeResponseItem
type ChangedFile struct {
	Kind uint8  `json:"kind"`
	Path string `json:"path"`
}
