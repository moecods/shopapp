package products

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(collection *mongo.Collection) *ProductRepository {
	return &ProductRepository{collection: collection}
}

func (r *ProductRepository) AddProduct(product *Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.collection.InsertOne(ctx, product)
	product.ID = result.InsertedID.(primitive.ObjectID)

	return err
}

func (r *ProductRepository) UpdateProduct(id primitive.ObjectID, product *Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": product,
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *ProductRepository) DeleteProduct(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *ProductRepository) ListProducts() ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(ctx)

	var products []Product
	for cursor.Next(ctx) {
		var product Product
		err := cursor.Decode(&product)
		if err != nil {
			log.Fatal(err)
		}
		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	return products, nil
}

func (r *ProductRepository) GetProduct(id primitive.ObjectID) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var product Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	return &product, err
}