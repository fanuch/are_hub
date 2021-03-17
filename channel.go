package are_server

import "context"

type Channel struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	common   `bson:"inline"`
}

type ChannelRepo interface {
	// Get all channels
	All(context.Context) ([]Channel, error)

	// Create a new channel
	Insert(context.Context, Archetype) error

	// Find a channel by its ID.
	FindID(context.Context, string) (*Channel, error)

	// Find and update a channel by its ID.
	UpdateID(context.Context, string, *Channel) error

	// Find and delete a channel by its ID.
	DeleteID(context.Context, string) (*Channel, error)

	// Get a count of channels.
	Count(context.Context) (int64, error)
}
