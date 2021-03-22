package are_server

import (
	"context"
	"fmt"
)

// password implements json.Marshaler.
type password string

// Prevent passwords from being marshaled and sent to clients by implementing
// json.Marshaler and returning a empty byte slice. This allows for passwords
// to be unmarshaled without issue.
func (p password) MarshalJSON() ([]byte, error) {
	return []byte(`""`), nil
}

// Represents a group of users listening to a data stream.
type Channel struct {
	Name     string `json:"name"`
	Password password
	common   `bson:"inline"`
}

// Create a new channel.
func NewChannel(name, pw string) *Channel {
	return &Channel{Name: name, Password: password(pw)}
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

// Retrieve the channel's password as a string.
func (c *Channel) PasswordStr() string {
	return string(c.Password)
}

// Mutate the channel's password to the string provided.
func (c *Channel) SetPasswordStr(pw string) {
	c.Password = password(pw)
}

type ChannelRepo interface {
	// Get all channels
	All(context.Context) ([]Channel, error)

	// Create a new channel
	Insert(context.Context, Archetype) error

	// Find a channel by its ID.
	FindID(context.Context, string) (*Channel, error)

	// Find and update a channel by its ID.
	UpdateID(context.Context, string, Archetype) error

	// Find and delete a channel by its ID.
	DeleteID(context.Context, string) (*Channel, error)

	// Get a count of channels.
	Count(context.Context) (int64, error)
}
