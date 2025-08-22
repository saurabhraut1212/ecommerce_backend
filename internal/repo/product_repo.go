package repo

import (
	"context"
	"time"

	"github.com/saurabhraut1212/ecommerce_backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepo struct {
	col *mongo.Collection
}

func NewProductRepo(db *mongo.Database) *ProductRepo {
	return &ProductRepo{
		col: db.Collection("products"),
	}
}

func (r *ProductRepo) Create(ctx context.Context, p *models.Product) error {
	p.ID = primitive.NewObjectID()
	now := time.Now().UTC()
	p.CreatedAt, p.UpdatedAt = now, now
	_, err := r.col.InsertOne(ctx, p)
	return err
}

func (r *ProductRepo) GetById(ctx context.Context, id primitive.ObjectID) (*models.Product, error) {
	var p models.Product
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&p)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &p, err
}

func (r *ProductRepo) List(ctx context.Context, page, limit int) ([]models.Product, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	skip := int64((page - 1) * limit)

	cur, err := r.col.Find(ctx, bson.M{}, &options.FindOptions{
		Skip:  &skip,
		Limit: func(i int64) *int64 { return &i }(int64(limit)),
		Sort:  bson.M{"created_at": -1},
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []models.Product
	for cur.Next(ctx) {
		var p models.Product
		if err := cur.Decode(&p); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, cur.Err()
}

func (r *ProductRepo) Update(ctx context.Context, id primitive.ObjectID, update bson.M) (*models.Product, error) {
	update["updated_at"] = time.Now().UTC()
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var p models.Product
	err := r.col.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": update}, opts).Decode(&p)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &p, err
}

func (r *ProductRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
