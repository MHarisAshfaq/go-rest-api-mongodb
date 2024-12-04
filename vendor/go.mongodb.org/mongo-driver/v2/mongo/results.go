// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/operation"
)

// BulkWriteResult is the result type returned by a BulkWrite operation.
type BulkWriteResult struct {
	// The number of documents inserted.
	InsertedCount int64

	// The number of documents matched by filters in update and replace operations.
	MatchedCount int64

	// The number of documents modified by update and replace operations.
	ModifiedCount int64

	// The number of documents deleted.
	DeletedCount int64

	// The number of documents upserted by update and replace operations.
	UpsertedCount int64

	// A map of operation index to the _id of each upserted document.
	UpsertedIDs map[int64]interface{}

	// Operation performed with an acknowledged write. Values for other fields may
	// not be deterministic if the write operation was unacknowledged.
	Acknowledged bool
}

// InsertOneResult is the result type returned by an InsertOne operation.
type InsertOneResult struct {
	// The _id of the inserted document. A value generated by the driver will be of type bson.ObjectID.
	InsertedID interface{}

	// Operation performed with an acknowledged write. Values for other fields may
	// not be deterministic if the write operation was unacknowledged.
	Acknowledged bool
}

// InsertManyResult is a result type returned by an InsertMany operation.
type InsertManyResult struct {
	// The _id values of the inserted documents. Values generated by the driver will be of type bson.ObjectID.
	InsertedIDs []interface{}

	// Operation performed with an acknowledged write. Values for other fields may
	// not be deterministic if the write operation was unacknowledged.
	Acknowledged bool
}

// TODO(GODRIVER-2367): Remove the BSON struct tags on DeleteResult.

// DeleteResult is the result type returned by DeleteOne and DeleteMany operations.
type DeleteResult struct {
	DeletedCount int64 // The number of documents deleted.

	// Operation performed with an acknowledged write. Values for other fields may
	// not be deterministic if the write operation was unacknowledged.
	Acknowledged bool
}

// RewrapManyDataKeyResult is the result of the bulk write operation used to update the key vault collection with
// rewrapped data keys.
type RewrapManyDataKeyResult struct {
	*BulkWriteResult
}

// ListDatabasesResult is a result of a ListDatabases operation.
type ListDatabasesResult struct {
	// A slice containing one DatabaseSpecification for each database matched by the operation's filter.
	Databases []DatabaseSpecification

	// The total size of the database files of the returned databases in bytes.
	// This will be the sum of the SizeOnDisk field for each specification in Databases.
	TotalSize int64
}

func newListDatabasesResultFromOperation(res operation.ListDatabasesResult) ListDatabasesResult {
	var ldr ListDatabasesResult
	ldr.Databases = make([]DatabaseSpecification, 0, len(res.Databases))
	for _, spec := range res.Databases {
		ldr.Databases = append(
			ldr.Databases,
			DatabaseSpecification{Name: spec.Name, SizeOnDisk: spec.SizeOnDisk, Empty: spec.Empty},
		)
	}
	ldr.TotalSize = res.TotalSize
	return ldr
}

// DatabaseSpecification contains information for a database. This type is returned as part of ListDatabasesResult.
type DatabaseSpecification struct {
	Name       string // The name of the database.
	SizeOnDisk int64  // The total size of the database files on disk in bytes.
	Empty      bool   // Specifies whether or not the database is empty.
}

// UpdateResult is the result type returned from UpdateOne, UpdateMany, and ReplaceOne operations.
type UpdateResult struct {
	MatchedCount  int64       // The number of documents matched by the filter.
	ModifiedCount int64       // The number of documents modified by the operation.
	UpsertedCount int64       // The number of documents upserted by the operation.
	UpsertedID    interface{} // The _id field of the upserted document, or nil if no upsert was done.

	// Operation performed with an acknowledged write. Values for other fields may
	// not be deterministic if the write operation was unacknowledged.
	Acknowledged bool
}

// IndexSpecification represents an index in a database. This type is returned by the IndexView.ListSpecifications
// function and is also used in the CollectionSpecification type.
type IndexSpecification struct {
	// The index name.
	Name string

	// The namespace for the index. This is a string in the format "databaseName.collectionName".
	Namespace string

	// The keys specification document for the index.
	KeysDocument bson.Raw

	// The index version.
	Version int32

	// The length of time, in seconds, for documents to remain in the collection. The default value is 0, which means
	// that documents will remain in the collection until they're explicitly deleted or the collection is dropped.
	ExpireAfterSeconds *int32

	// If true, the index will only reference documents that contain the fields specified in the index. The default is
	// false.
	Sparse *bool

	// If true, the collection will not accept insertion or update of documents where the index key value matches an
	// existing value in the index. The default is false.
	Unique *bool

	// The clustered index.
	Clustered *bool
}

type indexListSpecificationResponse struct {
	Name               string   `bson:"name"`
	Namespace          string   `bson:"ns"`
	KeysDocument       bson.Raw `bson:"key"`
	Version            int32    `bson:"v"`
	ExpireAfterSeconds *int32   `bson:"expireAfterSeconds"`
	Sparse             *bool    `bson:"sparse"`
	Unique             *bool    `bson:"unique"`
	Clustered          *bool    `bson:"clustered"`
}

// CollectionSpecification represents a collection in a database. This type is returned by the
// Database.ListCollectionSpecifications function.
type CollectionSpecification struct {
	// The collection name.
	Name string

	// The type of the collection. This will either be "collection" or "view".
	Type string

	// Whether or not the collection is readOnly. This will be false for MongoDB versions < 3.4.
	ReadOnly bool

	// The collection UUID. This field will be nil for MongoDB versions < 3.6. For versions 3.6 and higher, this will
	// be a bson.Binary with Subtype 4.
	UUID *bson.Binary

	// A document containing the options used to construct the collection.
	Options bson.Raw

	// An IndexSpecification instance with details about the collection's _id index.
	IDIndex IndexSpecification
}

// DistinctResult represents an array of BSON data returned from an operation.
// If the operation resulted in an error, all DistinctResult methods will return
// that error. If the operation did not return any data, all DistinctResult
// methods will return ErrNoDocuments.
type DistinctResult struct {
	err      error
	arr      bson.RawArray
	reg      *bson.Registry
	bsonOpts *options.BSONOptions
}

// Decode will unmarshal the array represented by this DistinctResult into v. If
// there was an error from the operation that created this DistinctReuslt, that
// error will be returned. If the operation returned no array, Decode will
// return ErrNoDocuments.
//
// If the operation was successful and returned an array, Decode will return any
// errors from the unmarshalling process without any modification. If v is nil
// or is a typed nil, an error will be returned.
func (dr *DistinctResult) Decode(v any) error {
	doc := bsoncore.NewDocumentBuilder().
		AppendValue("arr", bsoncore.Value{
			Type: bsoncore.TypeArray,
			Data: dr.arr,
		}).Build()

	dec := getDecoder(doc, dr.bsonOpts, dr.reg)

	return dec.Decode(&struct{ Arr any }{Arr: v})
}

// Err provides a way to check for query errors without calling Decode. Err
// returns the error, if any, that was encountered while running the operation.
// If the operation was successful but did not return any documents, Err returns
// ErrNoDocuments. If this error is not nil, this error will also be returned
// from Decode.
func (dr *DistinctResult) Err() error {
	return dr.err
}

// Raw returns the document represented by this DistinctResult as a bson.Raw. If
// there was an error from the operation that created this DistinctResult, both
// the result and that error will be returned. If the operation returned no
// documents, this will return (nil, ErrNoDocuments).
func (dr *DistinctResult) Raw() (bson.RawArray, error) {
	if dr.err != nil {
		return nil, dr.err
	}

	return dr.arr, nil
}