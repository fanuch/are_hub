package are_server

// All types should implement this interface so as to be compatible with the various
// sub-packages.
type Archetype interface {
	// Should mutate the ID field to the value of the paramter.
	SetID(string)

	// Should mutate the ID field to a zero value. Eg. the empty string ("")
	UnsetID()

	// Should mutate the CreatedAt and UpdatedAt (or alternatively named) fields
	// to the current time (preferably UTC+0)
	Created()

	// Should only mutate the UpdateAt (or similar) field to the current
	// time (preferably UTC+0).
	Updated()
}
