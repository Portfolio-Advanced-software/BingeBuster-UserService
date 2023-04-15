package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	collection *mongo.Collection
}

func NewUserService(db *mongo.Database) *UserService {
	return &UserService{
		collection: db.Collection("users"),
	}
}

func (s *UserService) Create(ctx context.Context, user *User) (interface{}, error) {
	result, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (s *UserService) Read(ctx context.Context, id primitive.ObjectID) (*User, error) {
	var user *User
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Update(ctx context.Context, id primitive.ObjectID, user *User) error {
	_, err := s.collection.ReplaceOne(ctx, bson.M{"_id": id}, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}
