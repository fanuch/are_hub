package validate

import (
	"net/http"

	"github.com/blacksfk/are_server"
	"github.com/go-playground/validator/v10"
)

type Channel struct {
	request
}

// Create a new channel validator which exports methods matching the
// microframework.Middleware signature.
func NewChannel() Channel {
	// dependency injection? never heard of it...
	return Channel{request{validator.New()}}
}

type channelStore struct {
	Name     string `validate:"required,alphanumunicode"`
	Password string `validate:"required"`
}

// Validate the request body with rules defined above. If successful,
// create an are_server.Channel and attach it to the request's context.
func (c Channel) ChannelStore(r *http.Request) error {
	temp := channelStore{}
	e := c.bodyStruct(r, &temp)

	if e != nil {
		return e
	}

	// create a channel (domain type) out of the validation object
	channel := &are_server.Channel{
		Name:     temp.Name,
		Password: temp.Password,
	}

	// insert the create channel into r's context
	*r = *r.WithContext(channel.ToCtx(r.Context()))

	return nil
}
