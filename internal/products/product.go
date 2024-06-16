package products

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
    ID          primitive.ObjectID `bson:"_id,omitempty"`
    Name        string             `bson:"name"`
    Description string             `bson:"description"`
    Price       float64            `bson:"price"`
    Stock       int                `bson:"stock"`
}