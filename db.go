package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func mongo_collection(name string) (*mongo.Collection, bool, error) {
	ctx, _ := context.WithTimeout(context.Background(), 20 * time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, false, err
	}
	collection := client.Database("Kotoba").Collection(name)
	return collection, true, nil
}

func count(name string, filter bson.M) (int, error) {
	collection, status, err := mongo_collection(name)
	if !status {
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

func get_one(name string, filter bson.M) (bson.M, error) {
	collection, status, err := mongo_collection(name)
	if !status {
		return nil, err
	}

	var result bson.M

	fndErr := collection.FindOne(context.TODO(), filter).Decode(&result)
	
	if fndErr != nil {
		return nil, fndErr
	}
	return result, nil
}

func upsert_one(name string, filter bson.M, update bson.M) (bool, error) {
	collection, status, err := mongo_collection(name)
	if !status {
		return false, err
	}
	
	_, upsErr := collection.UpdateOne(context.TODO(), filter, update, options.Update().SetUpsert(true))
	if upsErr != nil {
		return false, upsErr
	}

	return true, nil
}

func delete_many(name string, filter bson.M) (int, error) {
	collection, status, err := mongo_collection(name)
	if !status {
		return 0, err
	}
	res, delErr := collection.DeleteMany(context.TODO(), filter)
	if delErr != nil {
		return 0, delErr
	}
	return int(res.DeletedCount), nil
}