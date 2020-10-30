package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/griddis/atlant_test/tools/logging"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongodbStore struct {
	logger *logging.Logger
	client IMongoClient
	coll   IMongoCollection
}

func NewMongoStore(ctx context.Context, conn, database, coll string) (Repository, error) {
	logger := logging.FromContext(ctx)
	logger = logger.With("repo", "mongodb")
	client, err := mongo.NewClient(options.Client().ApplyURI(conn))
	if err != nil {
		return nil, errors.Wrap(err, "Mongodb NewClient error")
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Mongodb connect error")
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Mongodb ping error")
	}
	collection := client.Database(database).Collection(coll)
	return &mongodbStore{
		logger: logger,
		client: &MongoClient{client},
		coll:   &MongoCollection{collection},
	}, nil
}

func (r *mongodbStore) Close(ctx context.Context) error {
	r.client.Disconnect(ctx)
	return nil
}
func (r *mongodbStore) UpdatePrice(ctx context.Context, obj ProductPrice) error {
	//search previous price
	var searchDocument ProductPrice
	filter := bson.M{"name": obj.Name, "price": obj.Price}
	opts := options.FindOne()
	err := r.coll.FindOne(context.TODO(), filter, opts).Decode(&searchDocument)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return errors.New("error from find document")
		}

	}

	if searchDocument.Price != obj.Price {
		var updatedDocument ProductPrice
		opts := options.FindOneAndUpdate().SetUpsert(true)
		filter := bson.D{{"name", obj.Name}}
		update := bson.D{{"$set", bson.M{"name": obj.Name, "price": obj.Price, "date": time.Now()}}, {"$inc", bson.M{"counter": 1}}}

		err := r.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedDocument)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil
			}
			return Noresult
		}
		r.logger.Debug("status", "document updated", "msg", fmt.Sprintf("%+v", updatedDocument))
		return nil
	}
	r.logger.Debug("status", "document don`t update", "msg", fmt.Sprintf("name: %s price=%f newprice=%f", obj.Name, searchDocument.Price, obj.Price))
	return nil
}

func (r *mongodbStore) ListPrice(ctx context.Context, sorter map[string]int32, limiter Limiter) ([]*ProductPrice, error) {

	var sort bson.D = bson.D{}
	if len(sorter) > 0 {
		for key, value := range sorter {
			sort = append(sort, bson.E{key, value})
		}
	}

	opts := options.Find().SetSort(sort)
	if limiter.Limit > 0 {
		opts.SetLimit(limiter.Limit)
	}
	filter := bson.M{}
	if len(limiter.Offsetbyid) > 0 {
		objID, _ := primitive.ObjectIDFromHex(limiter.Offsetbyid)
		filter = bson.M{"_id": bson.M{"$gt": objID}}
	}
	cursor, err := r.coll.Find(context.TODO(), filter, opts)
	if err != nil {
		log.Fatal(err)
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	var products []*ProductPrice
	for _, result := range results {
		product := &ProductPrice{
			ID:      result["_id"].(primitive.ObjectID).Hex(),
			Name:    result["name"].(string),
			Price:   float32(result["price"].(float64)),
			Date:    result["date"].(primitive.DateTime).Time(),
			Counter: uint32(result["counter"].(int32)),
		}
		products = append(products, product)
	}
	return products, nil
}
