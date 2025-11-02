package entity

type ObjectMetadata struct {
	ID              int64  `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	BucketID        int64  `json:"bucketID,omitempty"`
	LocationGroupID int64  `json:"locationGroupID,omitempty"`
}
