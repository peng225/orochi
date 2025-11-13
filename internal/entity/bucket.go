package entity

type BucketStatus string

const (
	BucketStatusActive  BucketStatus = "active"
	BucketStatusDeleted BucketStatus = "deleted"
)

type Bucket struct {
	ID         int64        `json:"id,omitempty"`
	Name       string       `json:"name,omitempty"`
	ECConfigID int64        `json:"ecConfigID,omitempty"`
	Status     BucketStatus `json:"status,omitempty"`
}
