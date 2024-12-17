package mongo

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/errors"
	"time"

	"github.com/kwaaka-team/orders-core/core/externalapi/database/drivers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	connectionTimeout = 5 * time.Second
	ensureIdxTimeout  = 5 * time.Second

	authClientsCollectionName = "auth_clients"
)

type Mongo struct {
	connURL string
	dbName  string

	client *mongo.Client
	DB     *mongo.Database

	connectionTimeout time.Duration
	ensureIdxTimeout  time.Duration

	authClientRepo drivers.AuthClientRepository
}

func (m *Mongo) Name() string { return "Mongo" }

func (m *Mongo) Connect() error {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	defer cancel()

	m.client, err = mongo.Connect(ctx, options.Client().ApplyURI(m.connURL))
	if err != nil {
		return err
	}

	if err := m.Ping(); err != nil {
		return err
	}

	m.DB = m.client.Database(m.dbName)

	return nil
}

func (m *Mongo) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	defer cancel()

	return m.client.Ping(ctx, readpref.Primary())
}

func (m *Mongo) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

func (m *Mongo) AuthClientRepository() drivers.AuthClientRepository {
	if m.authClientRepo == nil {
		m.authClientRepo = NewAuthClientRepository(
			m.DB.Collection(authClientsCollectionName),
		)
	}
	return m.authClientRepo
}

func New(conf drivers.DataStoreConfig) (drivers.DataStore, error) {
	if conf.URL == "" {
		return nil, errors.ErrInvalidConfigStruct
	}

	if conf.DataBaseName == "" {
		return nil, errors.ErrInvalidConfigStruct
	}

	return &Mongo{
		connURL:           conf.URL,
		dbName:            conf.DataBaseName,
		connectionTimeout: connectionTimeout,
		ensureIdxTimeout:  ensureIdxTimeout,
	}, nil
}
