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

type OrderRepo struct {
	col *mongo.Collection
}

func NewOrderRepo(db *mongo.Database) *OrderRepo {
	return &OrderRepo{
		col: db.Collection("orders"),
	}
}

func (r *OrderRepo) Create(ctx context.Context, o *models.Order) error {
	o.ID = primitive.NewObjectID()
	now := time.Now().UTC()
	o.CreatedAt, o.UpdatedAt = now, now
	if o.Status == "" {
		o.Status = "pending"
	}
	_, err := r.col.InsertOne(ctx, o)
	return err
}

func (r *OrderRepo) GetById(ctx context.Context, id primitive.ObjectID) (*models.Order, error) {
	var o models.Order
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&o)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &o, err
}

func (r *OrderRepo) ListByUser(ctx context.Context, userId primitive.ObjectID, page, limit int) ([]models.Order, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	skip := int64((page - 1) * limit)

	cur, err := r.col.Find(ctx, bson.M{"user_id": userId}, &options.FindOptions{
		Skip:  &skip,
		Limit: func(i int64) *int64 { return &i }(int64(limit)),
		Sort:  bson.M{"created_at": -1},
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []models.Order
	for cur.Next(ctx) {
		var o models.Order
		if err := cur.Decode(&o); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, cur.Err()
}

func (r *OrderRepo) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) (*models.Order, error) {
	update := bson.M{"status": status, "updated_at": time.Now().UTC()}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var o models.Order
	err := r.col.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": update}, opts).Decode(&o)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &o, err
}

func (r *OrderRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
