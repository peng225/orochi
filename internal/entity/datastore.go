package entity

type DatastoreStatus string

const (
	DatastoreStatusActive DatastoreStatus = "active"
	DatastoreStatusDown   DatastoreStatus = "down"
)

// FIXME: is it OK to add omitempty for required fields?
type Datastore struct {
	ID      int64           `json:"id,omitempty"`
	BaseURL string          `json:"baseURL,omitempty"`
	Status  DatastoreStatus `json:"status,omitempty"`
}
