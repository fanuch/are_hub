package mock

import (
	"context"

	"github.com/blacksfk/are_server"
)

// Implements are_server.ChannelRepo.
type ChannelRepo struct {
	AllFunc   func(context.Context) ([]are_server.Channel, error)
	AllCalled bool

	InsertFunc   func(context.Context, are_server.Archetype) error
	InsertCalled bool

	FindIDFunc   func(context.Context, string) (*are_server.Channel, error)
	FindIDCalled bool

	UpdateIDFunc   func(context.Context, string, are_server.Archetype) error
	UpdateIDCalled bool

	DeleteIDFunc   func(context.Context, string) (*are_server.Channel, error)
	DeleteIDCalled bool

	CountFunc   func(context.Context) (int64, error)
	CountCalled bool
}

func (r *ChannelRepo) All(ctx context.Context) ([]are_server.Channel, error) {
	r.AllCalled = true

	return r.AllFunc(ctx)
}

func (r *ChannelRepo) Insert(ctx context.Context, archetype are_server.Archetype) error {
	r.InsertCalled = true

	return r.InsertFunc(ctx, archetype)
}

func (r *ChannelRepo) FindID(ctx context.Context, id string) (*are_server.Channel, error) {
	r.FindIDCalled = true

	return r.FindIDFunc(ctx, id)
}

func (r *ChannelRepo) UpdateID(ctx context.Context, id string, channel are_server.Archetype) error {
	r.UpdateIDCalled = true

	return r.UpdateIDFunc(ctx, id, channel)
}

func (r *ChannelRepo) DeleteID(ctx context.Context, id string) (*are_server.Channel, error) {
	r.DeleteIDCalled = true

	return r.DeleteIDFunc(ctx, id)
}

func (r *ChannelRepo) Count(ctx context.Context) (int64, error) {
	return 0, nil
}
