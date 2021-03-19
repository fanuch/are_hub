package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB connection parameters.
type Params struct {
	// Name of the mongodb user to authenticate with.
	User string

	// Password of the said user.
	Password string

	// Authentication mechanism. Eg. "SCRAM-SHA-256", "MONGODB-X509".
	Mechanism string

	// Address (host + port). Eg. "localhost:27017".
	Address string

	// Name of the database.
	Name string
}

// Connect to the mongodb instance provided in p. See Params for more information.
func Connect(ctx context.Context, p *Params) (*mongo.Client, error) {
	opts := options.Client()
	creds := options.Credential{
		Username:      p.User,
		Password:      p.Password,
		AuthMechanism: p.Mechanism,
		AuthSource:    p.Name,
	}

	opts.SetAuth(creds)
	opts.SetHosts([]string{p.Address})

	return mongo.Connect(ctx, opts)
}
