package repository

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/models"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	NewsletterCollection = "newsletter"
)

type Repository struct {
	collection *mongo.Collection
}

func NewNewsletterRepository(db *mongo.Database) *Repository {
	return &Repository{
		collection: db.Collection(NewsletterCollection),
	}
}

func (r *Repository) CreateNewsletter(ctx context.Context, newsletter models.Newsletter) error {
	if _, err := r.collection.InsertOne(ctx, newsletter); err != nil {
		return err
	}
	return nil
}
