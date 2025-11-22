package entity

type LocationGroupStatus string

const (
	LocationGroupStatusActive   LocationGroupStatus = "active"
	LocationGroupStatusDeleting LocationGroupStatus = "deleting"
)

type LocationGroup struct {
	ID         int64               `json:"id,omitempty"`
	Datastores []int64             `json:"datastores,omitempty"`
	ECConfigID int64               `json:"ecConfigID,omitempty"`
	Status     LocationGroupStatus `json:"status,omitempty"`
}
