package http

import (
	"fmt"
	"sync"

	"github.com/blacksfk/are_hub"
	"nhooyr.io/websocket"
)

const (
	// Maximum number of clients listening.
	MAX_SUBS = 10
)

// Wraps are_hub.Channel with a publisher client, a map of subscriber clients,
// and mutual exclusion locks.
type telemetryChannel struct {
	*are_hub.Channel

	pubMtx sync.Mutex
	pub    *websocket.Conn

	subMtx sync.Mutex
	subs   map[string]*client
}

func newTelemetryChannel(c *are_hub.Channel) *telemetryChannel {
	subs := make(map[string]*client, MAX_SUBS)

	return &telemetryChannel{c, sync.Mutex{}, nil, sync.Mutex{}, subs}
}

// Add a publisher.
func (c *telemetryChannel) setPub(pub *websocket.Conn) error {
	if c.pub != nil {
		return wsChannelFull()
	}

	c.pubMtx.Lock()
	defer c.pubMtx.Unlock()

	c.pub = pub

	return nil
}

// Remove the current publisher.
func (c *telemetryChannel) removePub() {
	c.pubMtx.Lock()
	defer c.pubMtx.Unlock()

	c.pub = nil
}

// Add a subscriber to the telemetryChannel.
func (c *telemetryChannel) addSub(sub *client) error {
	if len(c.subs) == MAX_SUBS {
		return wsChannelFull()
	}

	// lock the mutex to prevent changes to the map
	c.subMtx.Lock()
	defer c.subMtx.Unlock()

	// check if the subscriber ID already exists
	_, ok := c.subs[sub.id]

	if ok {
		// subscriber ID already exists
		return fmt.Errorf("Subscriber ID %s already exists.", sub.id)
	}

	// ID doesn't exist; add the subscriber
	c.subs[sub.id] = sub

	return nil
}

func (c *telemetryChannel) removeSub(sub *client) {
	c.subMtx.Lock()
	defer c.subMtx.Unlock()

	delete(c.subs, sub.id)
}

// Send a message to all subscribers on this channel.
func (c *telemetryChannel) broadcast(bytes []byte) {
	// prevent modification to the current subscribers while sending to the buffered channels
	c.subMtx.Lock()
	defer c.subMtx.Unlock()

	for id, sub := range c.subs {
		select {
		case sub.buffer <- bytes:
		default:
			// subscriber's buffer is full; drop them like a 10 tonne hammer
			go sub.drop(WS_ERROR_TIMEOUT, "Message buffer full")
			delete(c.subs, id)
		}
	}
}
