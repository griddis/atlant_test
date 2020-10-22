package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:generate echo generate mongodb mock
//go:generate mockgen -source=mongodb_abstract_layer.go -destination=mongodb_abstract_layer_mock.go -package=repository IMongoCollection

// mongo client
type IMongoClient interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Database(name string, opts ...*options.DatabaseOptions) IMongoDatabase
}

type MongoClient struct {
	Client *mongo.Client
}

func (c *MongoClient) Connect(ctx context.Context) error {
	return c.Connect(ctx)
}

func (c *MongoClient) Disconnect(ctx context.Context) error {
	return c.Disconnect(ctx)
}
func (c *MongoClient) Database(name string, opts ...*options.DatabaseOptions) IMongoDatabase {
	return &MongoDatabase{c.Client.Database(name, opts...)}
}

//mongo database
type IMongoDatabase interface {
	Collection(name string, opts ...*options.CollectionOptions) IMongoCollection
}

type MongoDatabase struct {
	Database *mongo.Database
}

func (d *MongoDatabase) Collection(name string, opts ...*options.CollectionOptions) IMongoCollection {
	return &MongoCollection{d.Database.Collection(name, opts...)}
}

// mongo collection
type IMongoDeleteResult interface{}
type IMongoInsertOneResult interface{}

type IMongoCollection interface {
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (IMongoCursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) IMongoSingleResult
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (IMongoInsertOneResult, error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (IMongoDeleteResult, error)
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) IMongoSingleResult
}

type IMongoSingleResult interface {
	Decode(v interface{}) error
}

type IMongoCursor interface {
	Close(context.Context) error
	Next(context.Context) bool
	Decode(interface{}) error
	All(context.Context, interface{}) error
}

type MongoCollection struct {
	Сoll *mongo.Collection
}

type MongoSingleResult struct {
	sr *mongo.SingleResult
}

type MongoCursor struct {
	cur *mongo.Cursor
}

func (msr *MongoSingleResult) Decode(v interface{}) error {
	return msr.sr.Decode(v)
}

func (mc *MongoCursor) Close(ctx context.Context) error {
	return mc.cur.Close(ctx)
}

func (mc *MongoCursor) Next(ctx context.Context) bool {
	return mc.cur.Next(ctx)
}

func (mc *MongoCursor) Decode(val interface{}) error {
	return mc.cur.Decode(val)
}

func (mc *MongoCursor) All(ctx context.Context, val interface{}) error {
	return mc.cur.All(ctx, val)
}

func (mc *MongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (IMongoCursor, error) {
	cursorResult, err := mc.Сoll.Find(ctx, filter, opts...)
	return &MongoCursor{cur: cursorResult}, err
}

func (mc *MongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) IMongoSingleResult {
	singleResult := mc.Сoll.FindOne(ctx, filter, opts...)
	return &MongoSingleResult{sr: singleResult}
}

func (mc *MongoCollection) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) IMongoSingleResult {
	singleResult := mc.Сoll.FindOneAndUpdate(ctx, filter, update, opts...)
	return &MongoSingleResult{sr: singleResult}
}

func (mc *MongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (IMongoInsertOneResult, error) {
	return mc.Сoll.InsertOne(ctx, document, opts...)
}

func (mc *MongoCollection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (IMongoDeleteResult, error) {
	return mc.Сoll.DeleteMany(ctx, filter, opts...)
}
