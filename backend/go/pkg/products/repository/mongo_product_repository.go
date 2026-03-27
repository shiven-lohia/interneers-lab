package repository

import (
	"context"
	"errors"
	"time"

	"github.com/shiven-lohia/interneers-lab/pkg/products/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoProductRepository struct {
	collection *mongo.Collection
}

func NewMongoProductRepository(client *mongo.Client) *MongoProductRepository {
	collection := client.Database("inventory").Collection("products")

	return &MongoProductRepository{
		collection: collection,
	}
}

func (r *MongoProductRepository) Create(ctx context.Context, product entity.Product) (entity.Product, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.M{"id": product.ID})
	if err != nil {
		return product, err
	}

	if count > 0 {
		return product, errors.New("product with this id already exists")
	}

	_, err = r.collection.InsertOne(ctx, product)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return product, errors.New("product with this id already exists")
		}

		return product, err
	}

	return product, nil
}

func (r *MongoProductRepository) GetAll(ctx context.Context, categoryID string) ([]entity.Product, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if categoryID != "" {
		filter["category_id"] = categoryID
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var products []entity.Product
	err = cursor.All(ctx, &products)
	if err != nil {
		return nil, err
	}

	if products == nil {
		products = []entity.Product{}
	}

	return products, err
}

func (r *MongoProductRepository) GetByID(ctx context.Context, id string) (entity.Product, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var product entity.Product

	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&product)

	return product, err
}

func (r *MongoProductRepository) Update(ctx context.Context, id string, product entity.Product) (entity.Product, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	product.ID = id

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"id": id},
		bson.M{"$set": product},
	)
	if err != nil {
		return product, err
	}

	if result.MatchedCount == 0 {
		return product, errors.New("Product not found")
	}

	return product, nil
}

func (r *MongoProductRepository) Delete(ctx context.Context, id string) error {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("Product not found")
	}

	return nil
}
