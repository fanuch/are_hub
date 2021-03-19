package mongodb

import (
	"context"

	"github.com/blacksfk/are_server"
	"go.mongodb.org/mongo-driver/mongo"
)

// Implements are_server.ChannelRepo.
type Channel struct {
	collection
}

// Create a channels collection in db.
func NewChannelCollection(client *mongo.Client, db string) Channel {
	return Channel{collection{client, db, "channels"}}
}

func (c Channel) All(ctx context.Context) ([]are_server.Channel, error) {
	var channels []are_server.Channel

	return channels, c.all(ctx, &channels)
}

func (c Channel) FindID(ctx context.Context, id string) (*are_server.Channel, error) {
	channel := &are_server.Channel{}

	return channel, c.findID(ctx, id, channel)
}

func (c Channel) DeleteID(ctx context.Context, id string) (*are_server.Channel, error) {
	channel := &are_server.Channel{}

	return channel, c.deleteID(ctx, id, channel)
}
