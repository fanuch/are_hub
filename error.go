package are_server

import "fmt"

// Returned from repositories when a query returns nothing.
// NoObjectsFound implements error.
type NoObjectsFound struct {
	repo, query string
}

// Create a new NoObjectsFoundError specifying the repository (eg. "users")
// and the query value (eg. "id == 7").
func NewNoObjectsFound(r, q string) *NoObjectsFound {
	return &NoObjectsFound{r, q}
}

func (e *NoObjectsFound) Error() string {
	return fmt.Sprintf("%s: No objects found matching: %s", e.repo, e.query)
}

// Auxiliary function to determine whether or not an error is a NoObjectsFound error.
func IsNoObjectsFound(e error) bool {
	_, ok := e.(*NoObjectsFound)

	return ok
}
