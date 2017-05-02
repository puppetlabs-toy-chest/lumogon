package types

// ContainerReport contains the gathered capabilities from a single container
type ContainerReport struct {
	Schema               string                `json:"$schema"`
	Generated            string                `json:"generated"`
	ContainerReportID    string                `json:"container_report_id"`
	ContainerID          string                `json:"container_id"`
	HarvesterContainerID string                `json:"-"`
	ContainerName        string                `json:"container_name"`
	Capabilities         map[string]Capability `json:"capabilities"`
}
