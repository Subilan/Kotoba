package kotoba

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func mongoCollection(name string) (*mongo.Collection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	collection := client.Database("Kotoba").Collection(name)
	return collection, nil
}

func mongoCount(name string, filter bson.M) (int, error) {
	collection, err := mongoCollection(name)
	if err != nil {
		return -1, err
	}
	res, fndErr := collection.Find(context.TODO(), filter)
	if fndErr != nil {
		return -1, fndErr
	}
	count := 0
	var doc bool
	for {
		doc = res.Next(context.TODO())
		if !doc {
			break
		} else {
			count++
		}
	}
	return count, nil
}

func mongoGetOne(name string, filter bson.M) (bson.M, error) {
	collection, err := mongoCollection(name)
	if err != nil {
		return nil, err
	}

	var result bson.M

	fndErr := collection.FindOne(context.TODO(), filter).Decode(&result)

	if fndErr != nil {
		return nil, fndErr
	}
	return result, nil
}

func mongoInsertOne(name string, doc bson.M) error {
	collection, err := mongoCollection(name)
	if err != nil {
		return err
	}

	_, insErr := collection.InsertOne(context.TODO(), doc)

	if insErr != nil {
		return insErr
	}

	return nil
}

func mongoUpdateOne(name string, filter bson.M, update bson.M) error {
	collection, err := mongoCollection(name)
	if err != nil {
		return err
	}

	_, updErr := collection.UpdateOne(context.TODO(), filter, update)
	if updErr != nil {
		return updErr
	}

	return nil
}

func mongoDeleteMany(name string, filter bson.M) (int, error) {
	collection, err := mongoCollection(name)
	if err != nil {
		return 0, err
	}
	res, delErr := collection.DeleteMany(context.TODO(), filter)
	if delErr != nil {
		return 0, delErr
	}
	return int(res.DeletedCount), nil
}
