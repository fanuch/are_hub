package are_server

import "time"

type Archetype interface {
	SetID(string)
	Created()
	Updated()
}

// Implements Archetype.
type common struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (c *common) SetID(id string) {
	c.ID = id
}

func (c *common) Created() {
	c.CreatedAt = time.Now().UTC()
	c.UpdatedAt = c.CreatedAt
}

func (c *common) Updated() {
	c.UpdatedAt = time.Now().UTC()
}
