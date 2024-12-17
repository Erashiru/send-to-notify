package repository

import (
	"context"
	"github.com/kwaaka-team/orders-core/service/error_solutions/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const errorSolutionCollectionName = "error_solutions"

type Repository interface {
	GetErrorSolutions(ctx context.Context) ([]models.ErrorSolution, error)
	GetErrorSolutionByCode(ctx context.Context, code string) (models.ErrorSolution, error)
	GetErrorSolutionByType(ctx context.Context, errorType string) (models.ErrorSolution, error)
	GetTimeoutErrSolutions(ctx context.Context) ([]models.ErrorSolution, error)
}

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) (*MongoRepository, error) {
	return &MongoRepository{collection: db.Collection(errorSolutionCollectionName)}, nil
}

func (r *MongoRepository) GetErrorSolutions(ctx context.Context) ([]models.ErrorSolution, error) {

	var errorSolutions []models.ErrorSolution

	cursor, err := r.collection.Find(ctx, bson.D{}, options.Find())
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var errorSolution models.ErrorSolution
		if err := cursor.Decode(&errorSolution); err != nil {
			return nil, err
		}
		errorSolutions = append(errorSolutions, errorSolution)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	cursor.Close(ctx)

	return errorSolutions, nil
}

func (r *MongoRepository) GetErrorSolutionByCode(ctx context.Context, code string) (models.ErrorSolution, error) {

	filter := bson.D{
		{Key: "code", Value: code},
	}

	var errSolution models.ErrorSolution

	if err := r.collection.FindOne(ctx, filter).Decode(&errSolution); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.ErrorSolution{}, nil
		}
		return models.ErrorSolution{}, err
	}

	return errSolution, nil
}

func (r *MongoRepository) GetErrorSolutionByType(ctx context.Context, errorType string) (models.ErrorSolution, error) {

	filter := bson.D{
		{Key: "type", Value: errorType},
	}

	var errSolution models.ErrorSolution

	if err := r.collection.FindOne(ctx, filter).Decode(&errSolution); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.ErrorSolution{}, nil
		}
		return models.ErrorSolution{}, err
	}

	return errSolution, nil
}

func (r *MongoRepository) GetTimeoutErrSolutions(ctx context.Context) ([]models.ErrorSolution, error) {

	filter := bson.D{
		{Key: "is_timeout", Value: true},
	}

	var errorSolutions []models.ErrorSolution

	cursor, err := r.collection.Find(ctx, filter, options.Find())
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var errorSolution models.ErrorSolution
		if err := cursor.Decode(&errorSolution); err != nil {
			return nil, err
		}
		errorSolutions = append(errorSolutions, errorSolution)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	cursor.Close(ctx)

	return errorSolutions, nil
}
