package dto

import "go.mongodb.org/mongo-driver/mongo"

type Config struct {
	Region    string
	SecretEnv string
	MongoCli  *mongo.Client
}
