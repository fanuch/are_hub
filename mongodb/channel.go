package mongodb

import (
	"context"

	"github.com/blacksfk/are_hub"
	"go.mongodb.org/mongo-driver/mongo"
)

// Implements are_hub.ChannelRepo.
type Channel struct {
	collection
}

// Create a channels collection in db.
func NewChannelCollection(client *mongo.Client, db string) Channel {
	return Channel{collection{client, db, "channels"}}
}

func (c Channel) All(ctx context.Context) ([]are_hub.Channel, error) {
	var channels []are_hub.Channel

	return channels, c.all(ctx, &channels)
}

func (c Channel) FindID(ctx context.Context, id string) (*are_hub.Channel, error) {
	channel := &are_hub.Channel{}

	return channel, c.findID(ctx, id, channel)
}

func (c Channel) DeleteID(ctx context.Context, id string) (*are_hub.Channel, error) {
	channel := &are_hub.Channel{}

	return channel, c.deleteID(ctx, id, channel)
}
