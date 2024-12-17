package repository

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/integration_api/models"
	"github.com/kwaaka-team/orders-core/core/integration_api/repository/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type DBInfo struct {
	Driver      string
	DBName      string
	MongoClient *mongo.Client
}

func NewExternalPosIntegrationAuthRepository(databaseInfo DBInfo) (ExternalPosIntegrationAuthRepository, error) {
	switch databaseInfo.Driver {
	case "mongo":
		return mongodb.NewExternalPosIntegrationAuthRepository(databaseInfo.MongoClient, databaseInfo.DBName), nil
	}

	return nil, fmt.Errorf("driver is not valid")
}

type ExternalPosIntegrationAuthRepository interface {
	GetAuthInfo(ctx context.Context, token string) (models.AuthInfo, error)
}
