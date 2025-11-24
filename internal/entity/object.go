package entity

type ObjectStatus string

const (
	ObjectStatusCreating ObjectStatus = "creating"
	ObjectStatusUpdating ObjectStatus = "updating"
	ObjectStatusActive   ObjectStatus = "active"
)

type ObjectMetadata struct {
	ID              int64        `json:"id,omitempty"`
	Name            string       `json:"name,omitempty"`
	Status          ObjectStatus `json:"status,omitempty"`
	BucketID        int64        `json:"bucketID,omitempty"`
	LocationGroupID int64        `json:"locationGroupID,omitempty"`
}
