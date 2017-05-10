package types

// TargetContainer contains the information needed to address a
// specific container whose data will be gathered
type TargetContainer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	OSID string `json:"os_id"`
}
