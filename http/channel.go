package http

import (
	"net/http"

	"github.com/blacksfk/are_server"
)

// CRUD controller that manipulates channel data in the provided repository.
type Channel struct {
	channels are_server.ChannelRepo
}

// Create a new channel controller.
func NewChannel(channels are_server.ChannelRepo) Channel {
	return Channel{channels}
}

// TODO: implement
func (c Channel) Index(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)

	return nil
}

// TODO: implement
func (c Channel) Store(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)

	return nil
}

// TODO: implement
func (c Channel) Show(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)

	return nil
}

// TODO: implement
func (c Channel) Update(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)

	return nil
}

// TODO: implement
func (c Channel) Delete(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)

	return nil
}
