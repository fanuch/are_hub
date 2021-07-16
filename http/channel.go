package http

import (
	"net/http"

	"github.com/blacksfk/are_hub"
	"github.com/blacksfk/are_hub/hash"
	uf "github.com/blacksfk/microframework"
)

// CRUD controller that manipulates channel data in the provided repository.
type Channel struct {
	channels are_hub.ChannelRepo
}

// Create a new channel controller.
func NewChannel(channels are_hub.ChannelRepo) Channel {
	return Channel{channels}
}

// Get all channels.
func (c Channel) Index(w http.ResponseWriter, r *http.Request) error {
	h := w.Header()

	h.Set("Access-Control-Allow-Methods", h.Get("Allow"))
	h.Set("Access-Control-Allow-Origin", "*")

	channels, e := c.channels.All(r.Context())

	if e != nil {
		return e
	}

	return uf.SendJSON(w, channels)
}

// Create a new channel.
func (c Channel) Store(w http.ResponseWriter, r *http.Request) error {
	// get the channel from the request's context
	channel, e := are_hub.ChannelFromCtx(r.Context())

	if e != nil {
		return e
	}

	// hash the new channel's password
	hash, e := hash.Password(channel.PasswordStr())

	if e != nil {
		return e
	}

	// replace the channel's plaintext password with the generated hash
	channel.SetPasswordStr(hash)

	// insert the new channel into the repository
	e = c.channels.Insert(r.Context(), channel)

	if e != nil {
		return e
	}

	// return the created channel with the ID and timestamps
	return uf.SendJSON(w, channel)
}

// Find a specific channel by its ID.
func (c Channel) Show(w http.ResponseWriter, r *http.Request) error {
	channel, e := c.channels.FindID(r.Context(), uf.GetParam(r, "id"))

	if e != nil {
		if are_hub.IsNoObjectsFound(e) {
			return uf.NotFound(e.Error())
		}

		return e
	}

	return uf.SendJSON(w, channel)
}

// Update a specific channel by its ID.
func (c Channel) Update(w http.ResponseWriter, r *http.Request) error {
	// get the channel embedded in the request's context
	channel, e := are_hub.ChannelFromCtx(r.Context())

	if e != nil {
		return e
	}

	// hash the new password
	hash, e := hash.Password(channel.PasswordStr())

	if e != nil {
		return e
	}

	// replace the channel's plaintext password with the generated hash
	channel.SetPasswordStr(hash)

	// update the channel in the repository
	e = c.channels.UpdateID(r.Context(), uf.GetParam(r, "id"), channel)

	if e != nil {
		if are_hub.IsNoObjectsFound(e) {
			return uf.NotFound(e.Error())
		}

		return e
	}

	// return the updated channel
	return uf.SendJSON(w, channel)
}

// Delete a specific channel by its ID.
func (c Channel) Delete(w http.ResponseWriter, r *http.Request) error {
	// delete and get the deleted channel
	channel, e := c.channels.DeleteID(r.Context(), uf.GetParam(r, "id"))

	if e != nil {
		if are_hub.IsNoObjectsFound(e) {
			return uf.NotFound(e.Error())
		}

		return e
	}

	// return the deleted channel
	return uf.SendJSON(w, channel)
}
