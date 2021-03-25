package http

import (
	"fmt"

	"nhooyr.io/websocket"
)

// Application-defined status codes.
// See: https://tools.ietf.org/html/rfc6455#section-7.4.2
//
// Challenge responses
const (
	// Challenge passed.
	WS_CHALLENGE_SUCCESS = iota + 4000

	// Server expects the next client message to contain a password.
	WS_CHALLENGE_PASSWORD
)

// Informational
const (
	// No error occurred in the last message.
	WS_OK = iota + 4200
)

// Error codes
const (
	// Malformed message data received.
	WS_ERROR_BAD_MSG = iota + 4400

	// Invalid login attempt.
	WS_ERROR_UNAUTHORISED

	// Permission denied.
	WS_ERROR_FORBIDDEN

	// Object not found.
	WS_ERROR_NOT_FOUND

	// Took too long responding.
	WS_ERROR_TIMEOUT

	// Publisher/subscriber channel is full and the client cannot be added
	// until a publisher/subscriber disconnects.
	WS_ERROR_CHANNEL_FULL
)

// Encapsulates data and challenge messages.
type response struct {
	Status websocket.StatusCode `json:"status"`
	Data   interface{}          `json:"data"`
}

func challengeSucceededResponse() response {
	return response{Status: WS_CHALLENGE_SUCCESS}
}

func passwordChallengeResponse() response {
	return response{Status: WS_CHALLENGE_PASSWORD}
}

func dataResponse(data interface{}) response {
	return response{Status: WS_OK, Data: data}
}

// Encapsulates error code messages.
type errorResponse struct {
	Status  websocket.StatusCode `json:"status"`
	Message string               `json:"message"`
}

// errorResponse implements error.
func (er *errorResponse) Error() string {
	return fmt.Sprintf("WebSocket error %d: %s", er.Status, er.Message)
}

func wsPolicyViolation(str string) *errorResponse {
	return &errorResponse{websocket.StatusPolicyViolation, str}
}

func wsBadMsg(str string) *errorResponse {
	return &errorResponse{WS_ERROR_BAD_MSG, str}
}

func wsUnauthorised(str string) *errorResponse {
	return &errorResponse{WS_ERROR_UNAUTHORISED, str}
}

func wsNotFound(str string) *errorResponse {
	return &errorResponse{WS_ERROR_NOT_FOUND, str}
}

func wsTimeout() *errorResponse {
	return &errorResponse{Status: WS_ERROR_TIMEOUT}
}

func wsChannelFull() *errorResponse {
	return &errorResponse{Status: WS_ERROR_CHANNEL_FULL}
}
