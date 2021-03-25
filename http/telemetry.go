package http

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/blacksfk/are_server"
	"github.com/blacksfk/are_server/hash"
	uf "github.com/blacksfk/microframework"
	"nhooyr.io/websocket"
)

var (
	// How long to wait in seconds until dropping the connection to a client.
	timeout time.Duration = time.Second * 10
)

// Handles receiving data from publishers and forwarding that data along to
// subscribers on the same channel.
type TelemetryServer struct {
	accept *websocket.AcceptOptions
	repo   are_server.ChannelRepo

	mtx      sync.Mutex
	channels map[string]*telemetryChannel
}

// Create a new telemetry server.
func NewTelemetryServer(o *websocket.AcceptOptions, r are_server.ChannelRepo) *TelemetryServer {
	return &TelemetryServer{
		accept:   o,
		repo:     r,
		channels: make(map[string]*telemetryChannel),
	}
}

// Handle upgrading publisher clients to a websocket.
func (ts *TelemetryServer) Publish(w http.ResponseWriter, r *http.Request) error {
	id := uf.GetParam(r, "id")

	// check that an ID was actually provided before upgrading the connection
	if len(id) == 0 {
		return uf.BadRequest("Expected channel ID as URL parameter")
	}

	conn, e := websocket.Accept(w, r, ts.accept)

	if e != nil {
		// the library writes its own HTTP error responses instead of letting the
		// user handle the HTTP errors themselves
		return nil
	}

	// connection established; HTTP handling has finished.
	go ts.publish(id, conn)

	return nil
}

// Handle upgrade subscriber clients to a websocket.
func (ts *TelemetryServer) Subscribe(w http.ResponseWriter, r *http.Request) error {
	id := uf.GetParam(r, "id")

	// check that an ID was actually provided before upgrading the connection
	if len(id) == 0 {
		return uf.BadRequest("Expected channel ID as URL parameter")
	}

	conn, e := websocket.Accept(w, r, ts.accept)

	if e != nil {
		// the library writes its own HTTP error responses instead of letting the
		// user handle the HTTP errors themselves
		return nil
	}

	// connection established; HTTP handling has finished.
	go ts.subscribe(id, conn)

	return nil
}

// Subscription procedure and message handling.
func (ts *TelemetryServer) subscribe(id string, conn *websocket.Conn) {
	tc, e := ts.procedure(id, conn)

	if e != nil {
		handleError(e, conn)

		return
	}

	// create a client out of the connection and add it to the channel
	sub := newClient(conn)
	e = tc.addSub(sub)

	if e != nil {
		handleError(e, conn)

		return
	}

	// send all good response
	bytes, e := json.Marshal(challengeSucceededResponse())

	if e != nil {
		handleError(e, conn)
		tc.removeSub(sub)

		return
	}

	e = writeTimeout(context.TODO(), timeout, conn, bytes)

	if e != nil {
		handleError(e, conn)
		tc.removeSub(sub)

		return
	}

	// handle sending/receiving of messages
	for {
		ctx := context.Background()

		select {
		case msg := <-sub.buffer:
			// message received, send it to the subscriber
			e = writeTimeout(ctx, timeout, conn, msg)

			if e != nil {
				handleError(e, conn)
				tc.removeSub(sub)

				return
			}
		case <-ctx.Done():
			// something went wrong
			handleError(ctx.Err(), conn)
			tc.removeSub(sub)

			return
		}
	}
}

// Publishing procedure and message handling.
func (ts *TelemetryServer) publish(id string, conn *websocket.Conn) {
	tc, e := ts.procedure(id, conn)

	if e != nil {
		handleError(e, conn)

		return
	}

	e = tc.setPub(conn)

	if e != nil {
		handleError(e, conn)

		return
	}

	// send all good response
	bytes, e := json.Marshal(challengeSucceededResponse())

	if e != nil {
		handleError(e, conn)
		tc.removePub()

		return
	}

	e = writeTimeout(context.Background(), timeout, conn, bytes)

	if e != nil {
		handleError(e, conn)
		tc.removePub()

		return
	}

	// handle sending/receiving of messages
	for {
		// if the publisher takes longer than a minute, disconnect them
		bytes, e := readTimeout(context.Background(), time.Minute, conn)

		if e != nil {
			// something broke; disconnect the publisher
			handleError(e, conn)
			tc.removePub()

			return
		}

		// send the data to all of the clients on the same telemetry channel
		tc.broadcast(bytes)
	}
}

// Procedure (protocol) that the connecting client is expected to follow to establish
// itself as a publisher/subscriber of the channel it requested. See the protocol documentation
// in the docs directory.
func (ts *TelemetryServer) procedure(id string, conn *websocket.Conn) (*telemetryChannel, error) {
	// find the channel based on the provided ID
	channel, e := ts.repo.FindID(context.TODO(), id)

	if e != nil {
		if are_server.IsNoObjectsFound(e) {
			return nil, wsNotFound(e.Error())
		}

		return nil, e
	}

	// channel found, ask for the password
	bytes, e := passwordChallenge(context.TODO(), conn)

	if e != nil {
		return nil, e
	}

	// compare the received password with the known password
	match, e := hash.CmpPassword(channel.PasswordStr(), string(bytes))

	if e != nil {
		// something went wrong with comparison
		return nil, e
	}

	if !match {
		return nil, wsUnauthorised("Incorrect password")
	}

	// check if the telemetry channel already exists in the map
	tc, ok := ts.channels[id]

	if !ok {
		// create a new telemetry channel out of the found one and add it to the map
		tc = newTelemetryChannel(channel)
		ts.addChannel(tc)
	}

	return tc, nil
}

func (ts *TelemetryServer) addChannel(channel *telemetryChannel) {
	ts.mtx.Lock()
	ts.channels[channel.ID] = channel
	ts.mtx.Unlock()
}

func passwordChallenge(ctx context.Context, conn *websocket.Conn) ([]byte, error) {
	bytes, e := json.Marshal(passwordChallengeResponse())

	if e != nil {
		return nil, e
	}

	// send the challenge to the client
	e = writeTimeout(ctx, timeout, conn, bytes)

	if e != nil {
		return nil, e
	}

	// receive the challenge response from the client
	return readTimeout(ctx, timeout, conn)
}

func writeTimeout(ctx context.Context, timeout time.Duration, conn *websocket.Conn, bytes []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return conn.Write(ctx, websocket.MessageText, bytes)
}

func readTimeout(ctx context.Context, timeout time.Duration, conn *websocket.Conn) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	mtype, bytes, e := conn.Read(ctx)

	if e != nil {
		// something went wrong while reading
		return nil, e
	}

	if mtype != websocket.MessageText {
		// only text data is supported
		return nil, wsPolicyViolation("Binary data is not supported")
	}

	return bytes, nil
}

func handleError(e error, conn *websocket.Conn) {
	if er, ok := e.(*errorResponse); ok {
		conn.Close(er.Status, er.Message)
	} else {
		conn.Close(websocket.StatusInternalError, e.Error())
	}
}
