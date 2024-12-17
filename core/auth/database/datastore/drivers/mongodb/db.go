package mongodb

import (
	"context"
	drivers2 "github.com/kwaaka-team/orders-core/core/auth/database/datastore/drivers"
	"github.com/kwaaka-team/orders-core/core/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

const (
	connectionTimeout = 3 * time.Second
	ensureIdxTimeout  = 300 * time.Second

	usersCollectionName = "users"
)

type Mongo struct {
	connURL string
	dbName  string

	client *mongo.Client
	DB     *mongo.Database

	connectionTimeout time.Duration
	ensureIdxTimeout  time.Duration

	userRepo drivers2.UserRepository
}

func (m *Mongo) Name() string { return "Mongo" }

func (m *Mongo) Connect() error {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	defer cancel()

	m.client, err = mongo.Connect(ctx, options.Client().ApplyURI(m.connURL)) //"mongodb://localhost:27016"
	if err != nil {
		return err
	}

	if err := m.Ping(); err != nil {
		return err
	}

	m.DB = m.client.Database(m.dbName)

	return m.ensureIndexes()
}

func (m *Mongo) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	defer cancel()

	return m.client.Ping(ctx, readpref.Primary())
}

func (m *Mongo) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

func (m *Mongo) AuthRepository() drivers2.UserRepository {
	if m.userRepo == nil {
		m.userRepo = NewUserRepository(m.DB.Collection(usersCollectionName))
	}
	return m.userRepo
}

func New(conf drivers2.DataStoreConfig) (drivers2.DataStore, error) {
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

func (m *Mongo) ensureIndexes() error {
	ctx := context.Background()
	if err := m.ensureUserIndexes(ctx); err != nil {
		return err
	}

	return nil
}

func (m *Mongo) ensureUserIndexes(ctx context.Context) (err error) {
	col := m.DB.Collection(usersCollectionName)

	existingIndexes, err := m.existingIndexes(ctx, col)
	if err != nil {
		return err
	}

	indexesMap := map[string]mongo.IndexModel{}

	indexes := make([]mongo.IndexModel, 0, len(indexesMap))

	for name, idx := range indexesMap {
		if _, ok := existingIndexes[name]; ok {
			continue
		}

		idx.Options.SetName(name)
		indexes = append(indexes, idx)
	}

	if len(indexes) == 0 {
		return nil
	}

	opts := options.CreateIndexes().SetMaxTime(m.ensureIdxTimeout)
	_, err = col.Indexes().CreateMany(ctx, indexes, opts)

	return err
}

func (m *Mongo) existingIndexes(ctx context.Context, collection *mongo.Collection) (map[string]struct{}, error) {
	cur, err := collection.Indexes().List(ctx)
	if err != nil {
		return nil, err
	}
	defer errors.CloseCur(cur)

	res := make(map[string]struct{})

	for cur.Next(ctx) {
		res[cur.Current.Lookup("name").StringValue()] = struct{}{}
	}

	if err = cur.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
