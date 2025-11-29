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
	BucketName      string       `json:"bucketName,omitempty"`
	LocationGroupID int64        `json:"locationGroupID,omitempty"`
}
