package Repos

import (
	"InternBorobitApp/Interfaces"
	"InternBorobitApp/Models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type MongoGameRepository struct {
	collection *mongo.Collection
}

func NewMongoGameRepository(client *mongo.Client, dbName, collectionName string) Interfaces.GameRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &MongoGameRepository{collection: collection}
}

func (repo *MongoGameRepository) Create(game *Models.Game) error {
	game.ID = primitive.NewObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := repo.collection.InsertOne(ctx, game)
	return err
}

func (repo *MongoGameRepository) GetByID(id string) (*Models.Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var game Models.Game
	err = repo.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&game)
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (repo *MongoGameRepository) Update(game *Models.Game) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(game.ID.Hex())
	if err != nil {
		return err
	}
	_, err = repo.collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": game})
	return err
}

func (repo *MongoGameRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = repo.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (repo *MongoGameRepository) List() ([]Models.Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := repo.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)

	var games []Models.Game
	for cursor.Next(ctx) {
		var game Models.Game
		if err = cursor.Decode(&game); err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return games, nil
}
