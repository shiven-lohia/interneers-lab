package repository

import (
	"context"
	"errors"
	"time"

	"github.com/shiven-lohia/interneers-lab/pkg/products/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoCategoryRepository struct {
	collection *mongo.Collection
}

func NewMongoCategoryRepository(col *mongo.Collection) *MongoCategoryRepository {
	return &MongoCategoryRepository{collection: col}
}

// expose collection for testing purposes
func (r *MongoCategoryRepository) Collection() *mongo.Collection {
	return r.collection
}

func (r *MongoCategoryRepository) GetByID(ctx context.Context, id string) (entity.ProductCategory, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var category entity.ProductCategory
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entity.ProductCategory{}, errors.New("category not found")
		}
		return entity.ProductCategory{}, err
	}

	return category, nil
}

func (r *MongoCategoryRepository) Create(
	ctx context.Context,
	category entity.ProductCategory,
) (entity.ProductCategory, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if category.ID == "" {
		category.ID = primitive.NewObjectID().Hex()
	}

	count, err := r.collection.CountDocuments(
		ctx,
		bson.M{"_id": category.ID},
	)
	if err != nil {
		return category, err
	}

	if count > 0 {
		return category, errors.New("category with this id already exists")
	}

	_, err = r.collection.InsertOne(ctx, category)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return category, errors.New("category with this id already exists")
		}
		return category, err
	}

	return category, nil
}

func (r *MongoCategoryRepository) GetAll(ctx context.Context) ([]entity.ProductCategory, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var categories []entity.ProductCategory
	err = cursor.All(ctx, &categories)
	if err != nil {
		return nil, err
	}

	if categories == nil {
		categories = []entity.ProductCategory{}
	}

	return categories, nil
}

func (r *MongoCategoryRepository) Update(ctx context.Context, id string, category entity.ProductCategory) (entity.ProductCategory, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	category.ID = id

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": category},
	)
	if err != nil {
		return category, err
	}

	if result.MatchedCount == 0 {
		return category, errors.New("category not found")
	}

	return category, nil
}

func (r *MongoCategoryRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("category not found")
	}

	return nil
}
