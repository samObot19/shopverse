package repository

import (
    "context"
    "errors"
    "github.com/samObot19/shopverse/product-service/models"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

type MongoProductRepository struct {
    collection *mongo.Collection
}

func NewMongoProductRepository(collection *mongo.Collection) *MongoProductRepository {
    return &MongoProductRepository{collection: collection}
}

func (r *MongoProductRepository) CreateProduct(ctx context.Context, product *models.Product) error {
    _, err := r.collection.InsertOne(ctx, product)
    return err
}

func (r *MongoProductRepository) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
    var product models.Product
    err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, errors.New("product not found")
        }
        return nil, err
    }
    return &product, nil
}

func (r *MongoProductRepository) GetAllProducts(ctx context.Context, filters map[string]interface{}) ([]*models.Product, error) {
    var products []*models.Product
    cursor, err := r.collection.Find(ctx, filters)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var product models.Product
        if err := cursor.Decode(&product); err != nil {
            return nil, err
        }
        products = append(products, &product)
    }
    return products, nil
}

func (r *MongoProductRepository) UpdateProduct(ctx context.Context, id string, updatedProduct *models.Product) error {
    _, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updatedProduct})
    return err
}

func (r *MongoProductRepository) DeleteProduct(ctx context.Context, id string) error {
    _, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
    return err
}

func (r *MongoProductRepository) UpdateStock(ctx context.Context, id string, quantity int) error {
    _, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"stock": quantity}})
    return err
}

func (r *MongoProductRepository) GetProductsByCategory(ctx context.Context, category string) ([]*models.Product, error) {
    return r.GetAllProducts(ctx, bson.M{"category": category})
}

func (r *MongoProductRepository) SearchProducts(ctx context.Context, query string) ([]*models.Product, error) {
    filter := bson.M{"$text": bson.M{"$search": query}}
    return r.GetAllProducts(ctx, filter)
}


