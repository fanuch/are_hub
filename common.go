package are_server

import "time"

// Implements Archetype and is intended to be used in domain types
// via struct composition.
type common struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Set ID.
func (c *common) SetID(id string) {
	c.ID = id
}

// Reset ID to the empty string ("").
func (c *common) UnsetID() {
	c.ID = ""
}

// Set CreatedAt and UpdateAt to the current time (UTC+0).
func (c *common) Created() {
	c.CreatedAt = time.Now().UTC()
	c.UpdatedAt = c.CreatedAt
}

// Set UpdatedAt to the current time (UTC+0).
func (c *common) Updated() {
	c.UpdatedAt = time.Now().UTC()
}
