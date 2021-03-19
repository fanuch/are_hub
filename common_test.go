package are_server

import (
	"testing"
	"time"
)

func TestSetID(t *testing.T) {
	expected := "abc123"
	c := &common{}

	c.SetID(expected)

	if c.ID != expected {
		t.Fatalf("Expected: %s. Actual: %s", expected, c.ID)
	}
}

func TestUnSetID(t *testing.T) {
	c := &common{ID: "abc123"}

	c.UnsetID()

	if c.ID != "" {
		t.Fatalf("Expected: \"\". Actual: %s", c.ID)
	}
}

func TestCreated(t *testing.T) {
	c := &common{}
	zero := time.Time{}

	c.Created()

	if c.CreatedAt == zero || c.UpdatedAt == zero {
		t.Fatalf("Expected: non-zero timestamps. Actual: CreatedAt: %v, UpdatedAt: %v.",
			c.CreatedAt, c.UpdatedAt)
	}
}

func TestUpdated(t *testing.T) {
	c := &common{}
	zero := time.Time{}

	c.Updated()

	if c.UpdatedAt == zero {
		t.Fatalf("Expected non-zero timestamp. Actual: %+v", c.UpdatedAt)
	}
}
