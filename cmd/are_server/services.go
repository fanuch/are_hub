package main

import (
	"context"
	"log"

	"github.com/blacksfk/are_server"
	"github.com/blacksfk/are_server/mongodb"
)

// Wraps various database tables and services that
// require initilisation.
type services struct {
	channels are_server.ChannelRepo
}

// Initialise various services and create mongodb collections based on conf. This function
// is only intended to be called from the main function therefore dies if it encounters
// an error creating a mongo.Client.
func initServices(conf *config) *services {
	client, e := mongodb.Connect(context.Background(), conf.MongoDB)

	if e != nil {
		// connection failed so die
		log.Fatal(e)
	}

	return &services{
		mongodb.NewChannelCollection(client, conf.MongoDB.Name),
	}
}
