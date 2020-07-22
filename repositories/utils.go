package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func isCollectionExists(ctx context.Context, db *mongo.Database, collectionName string) (bool, error) {
	isExists := false
	names, err := db.ListCollectionNames(ctx, bson.D{{"name", collectionName}})
	if err != nil {
		return false, err
	}

	for _, name := range names {
		if name == collectionName {
			isExists = true
			break
		}
	}
	return isExists, nil
}
