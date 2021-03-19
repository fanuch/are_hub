package validate

import (
	"fmt"
	"net/http"
	"strings"

	uf "github.com/blacksfk/microframework"
	"github.com/go-playground/validator/v10"
)

type request struct {
	validate *validator.Validate
}

// Validate the request body as the passed in struct pointer.
func (r request) bodyStruct(req *http.Request, ptr interface{}) error {
	e := uf.DecodeBodyJSON(req, ptr)

	if e != nil {
		return e
	}

	e = c.validate.Struct(ptr)

	if e != nil {
		return assertError(e)
	}

	// body passed validation
	return nil
}

// Attempts to assert the provided error as validator.ValidationErrors and create a
// microframework.HttpError (404 Not Found) with the offending field(s).
//
// If unsuccessful, assertion of validator.InvalidValidationError is attempted and
// a 500 Internal Server Error is created.
//
// If unsucessful, the error is simply returned.
func assertError(e error) error {
	ve, ok := e.(validator.ValidationErrors)

	if !ok {
		// programmer error whoops
		ive, ok := e.(validator.InvalidValidationError)

		if !ok {
			// something's really busted
			return e
		}

		return uf.InternalServerError(e.Error())
	}

	// asserted as user error so loop and build user error string
	b := strings.Builder{}

	for _, fe := range ve {
		b.WriteString(fe.StructNamespace())
		b.WriteString(": expected: ")
		b.WriteString(fe.ActualTag())
		b.WriteString(", received: ")

		// leave a space for the next error
		b.WriteString(fmt.Sprintf("%v. ", fe.Value()))
	}

	return uf.BadRequest(b.String())
}
