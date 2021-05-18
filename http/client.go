package http

import (
	"encoding/hex"
	"math/rand"
	"time"

	"nhooyr.io/websocket"
)

const (
	// How long (in bytes) the random IDs should be
	ID_LEN = 6

	// How many messages to queue for the client before the connection
	// is dropped for being too slow.
	BUF_LEN = 16
)

// Wraps a websocket.Conn and implements a message buffer.
type client struct {
	id     string
	conn   *websocket.Conn
	buffer chan []byte
}

// Create a new client.
func newClient(conn *websocket.Conn) *client {
	// allocate a byte slice
	bytes := make([]byte, ID_LEN)

	// firstly seed the default math/rand source with the current epoch timestamp.
	rand.Seed(time.Now().Unix())

	// math/rand.Read always returns a nil error.
	rand.Read(bytes)

	// encode the random bytes as hex
	id := hex.EncodeToString(bytes)

	return &client{id, conn, make(chan []byte, BUF_LEN)}
}

// Disconnect the client.
func (c *client) drop(code websocket.StatusCode, reason string) {
	c.conn.Close(code, reason)
}
