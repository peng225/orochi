package entity

type BucketStatus string

const (
	BucketStatusActive   BucketStatus = "active"
	BucketStatusDeleting BucketStatus = "deleting"
)

type Bucket struct {
	ID       int64        `json:"id,omitempty"`
	Name     string       `json:"name,omitempty"`
	ECConfig string       `json:"ecConfig,omitempty"`
	Status   BucketStatus `json:"status,omitempty"`
}
