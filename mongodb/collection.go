package mongodb

import (
	"context"
	"fmt"

	"github.com/blacksfk/are_hub"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Wrapper around mongo.Client and mongo.Collection. To be used in struct
// composition with collections.
type collection struct {
	// mongodb connection pool.
	client *mongo.Client

	// Name of the database in the mongodb instance.
	db string

	// Name of the collection in the database.
	name string
}

// Get the collection instance from the client connection pool.
func (c collection) get() *mongo.Collection {
	return c.client.Database(c.db).Collection(c.name)
}

// Get a count of the number of documents in a collection.
func (c collection) Count(ctx context.Context) (int64, error) {
	return c.get().EstimatedDocumentCount(ctx)
}

// Get all documents.
func (c collection) all(ctx context.Context, slice interface{}) error {
	cursor, e := c.get().Find(ctx, bson.D{})

	if e != nil {
		return e
	}

	return cursor.All(ctx, slice)
}

// Get a document matching the hexadecimal-encoded ID.
func (c collection) findID(ctx context.Context, hex string, ptr are_hub.Archetype) error {
	id, e := primitive.ObjectIDFromHex(hex)

	if e != nil {
		return e
	}

	coll := c.get()
	result := coll.FindOne(ctx, bson.M{"_id": id})

	if e := result.Err(); e != nil {
		if e == mongo.ErrNoDocuments {
			return are_hub.NewNoObjectsFound("id: "+hex, coll.Name())
		}

		return e
	}

	return c.get().FindOne(ctx, bson.M{"_id": id}).Decode(ptr)
}

// Insert a document.
func (c collection) Insert(ctx context.Context, ptr are_hub.Archetype) error {
	ptr.Created()

	result, e := c.get().InsertOne(ctx, ptr)

	if e != nil {
		return e
	}

	id, ok := result.InsertedID.(primitive.ObjectID)

	if !ok {
		// this should never happen, right?
		e = fmt.Errorf("Could not assert %v as primitive.ObjectID", result.InsertedID)

		return e
	}

	// set the generated ID
	ptr.SetID(id.Hex())

	return nil
}

// Update a document matching the hexadecimal-encoded ID.
func (c collection) UpdateID(ctx context.Context, hex string, ptr are_hub.Archetype) error {
	id, e := primitive.ObjectIDFromHex(hex)

	if e != nil {
		return e
	}

	// set the updated timestamp and unset the ID to prevent it from being overwritten
	ptr.Updated()
	ptr.UnsetID()

	// ensure the new document is returned rather than the document before updating
	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(options.After)

	return c.get().FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": ptr}, opts).Decode(ptr)
}

// Delete a document matching the hexadecimal-encoded ID.
func (c collection) deleteID(ctx context.Context, hex string, ptr are_hub.Archetype) error {
	id, e := primitive.ObjectIDFromHex(hex)

	if e != nil {
		return e
	}

	return c.get().FindOneAndDelete(ctx, bson.M{"_id": id}).Decode(ptr)
}
