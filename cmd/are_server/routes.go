package main

import (
	"github.com/blacksfk/are_server/http"
	"github.com/blacksfk/are_server/http/middleware/validate"
	uf "github.com/blacksfk/microframework"
)

// HTTP route definitions.
func routes(s *uf.Server, services *services) {
	// channel routes
	c := http.NewChannel(services.channels)
	v := validate.NewChannel()

	s.NewGroup("/channel").Get(c.Index).Post(c.Store, v.Store)
	s.NewGroup("/channel/:id").Get(c.Show).Put(c.Update, v.Store).Delete(c.Delete)
}
