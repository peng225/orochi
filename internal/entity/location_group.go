package entity

type LocationGroup struct {
	ID                int64   `json:"id,omitempty"`
	CurrentDatastores []int64 `json:"currentDatastores,omitempty"`
	DesiredDatastores []int64 `json:"desiredDatastores,omitempty"`
}
