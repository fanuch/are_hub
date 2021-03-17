package are_server

import (
	"context"
	"fmt"
)

// Represents a group of users listening to a data stream.
type Channel struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	common   `bson:"inline"`
}

// Get a channel from a context.
func ChannelFromCtx(ctx context.Context) (*Channel, error) {
	v := ctx.Value(keyChannel)
	c, ok := v.(*Channel)

	if !ok {
		return nil, fmt.Errorf("Could not assert %v as *Channel\n", v)
	}

	return c, nil
}

// Insert a channel into context.
func (c *Channel) ToCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, keyChannel, c)
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
