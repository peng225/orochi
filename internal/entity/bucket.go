package entity

type Bucket struct {
	ID         int64  `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	ECConfigID int64  `json:"ecConfigID,omitempty"`
	Status     string `json:"status,omitempty"`
}
