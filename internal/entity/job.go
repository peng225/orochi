package entity

const (
	DeleteAllObjectsInBucket = "DeleteAllObjectsInBucket"
)

type Job struct {
	ID   int64  `json:"id,omitempty"`
	Kind string `json:"kind,omitempty"`
	Data []byte `json:"data,omitempty"`
}

type DeleteAllObjectsInBucketParam struct {
	BucketID int64 `json:"bucketID,omitempty"`
}
