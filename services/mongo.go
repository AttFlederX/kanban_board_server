package services

import (
	"context"
	"time"

	"github.com/AttFlederX/kanban_board_server/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoService struct {
	CollectionName string
}

func NewMongoService(collectionName string) *MongoService {
	return &MongoService{
		CollectionName: collectionName,
	}
}

func (s *MongoService) FindAll(result any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.DB.Collection(s.CollectionName).Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, result)
}

func (s *MongoService) Find(filter bson.M, result any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.DB.Collection(s.CollectionName).Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, result)
}

func (s *MongoService) FindByID(id primitive.ObjectID, result any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return database.DB.Collection(s.CollectionName).FindOne(ctx, bson.M{"_id": id}).Decode(result)
}

func (s *MongoService) FindOne(filter bson.M, result any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return database.DB.Collection(s.CollectionName).FindOne(ctx, filter).Decode(result)
}

func (s *MongoService) InsertOne(document any) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := database.DB.Collection(s.CollectionName).InsertOne(ctx, document)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (s *MongoService) UpdateByID(id primitive.ObjectID, update bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := database.DB.Collection(s.CollectionName).UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (s *MongoService) DeleteByID(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := database.DB.Collection(s.CollectionName).DeleteOne(ctx, bson.M{"_id": id})
	return err
}
