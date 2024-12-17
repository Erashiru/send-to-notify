package mongo

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/errors"
	"time"

	"github.com/kwaaka-team/orders-core/core/database/drivers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	connectionTimeout = 3 * time.Second
	ensureIdxTimeout  = 300 * time.Second

	menuCollectionName        = "menus"
	menuUploadTransactionName = "menu_upload_transactions"
	orderCollectionName       = "orders"
	bkOfferCollectionName     = "bk_offers"
)

type Mongo struct {
	connURL string
	dbName  string

	client *mongo.Client
	DB     *mongo.Database

	connectionTimeout time.Duration
	ensureIdxTimeout  time.Duration

	orderRepo     drivers.OrderRepository
	bkOfferRepo   drivers.BKOfferRepository
	analyticsRepo drivers.AnalyticsRepository
}

func (m *Mongo) Name() string { return "Mongo" }

func (m *Mongo) Client() *mongo.Client {
	return m.client
}

func (m *Mongo) Connect(cli *mongo.Client) error {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	defer cancel()

	if err = m.initializeClient(ctx, cli); err != nil {
		return err
	}

	if err := m.Ping(); err != nil {
		return err
	}

	m.DB = m.client.Database(m.dbName)

	return m.ensureIndexes()
}

func (m *Mongo) initializeClient(ctx context.Context, cli *mongo.Client) error {
	if cli != nil {
		m.client = cli
		return nil
	}

	if m.client != nil {
		return nil
	}

	if client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.connURL)); err != nil {
		return err
	} else {
		m.client = client
	}
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

func (m *Mongo) OrderRepository() drivers.OrderRepository {
	if m.orderRepo == nil {
		m.orderRepo = NewOrderRepository(m.DB.Collection(orderCollectionName))
	}
	return m.orderRepo
}

func (m *Mongo) BKOfferRepository() drivers.BKOfferRepository {
	if m.bkOfferRepo == nil {
		m.bkOfferRepo = NewBKOfferRepository(m.DB.Collection(bkOfferCollectionName))
	}
	return m.bkOfferRepo
}

func (m *Mongo) AnalyticsRepository() drivers.AnalyticsRepository {
	if m.analyticsRepo == nil {
		m.analyticsRepo = NewAnalyticsRepository(m.DB.Collection(orderCollectionName))
	}
	return m.analyticsRepo
}

func New(conf drivers.DataStoreConfig) (*Mongo, error) {
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
	if err := m.ensureMenuIndexes(ctx); err != nil {
		return err
	}
	if err := m.ensureOrderIndexes(ctx); err != nil {
		return err
	}

	return nil
}

func (m *Mongo) ensureMenuIndexes(ctx context.Context) (err error) {
	col := m.DB.Collection(menuCollectionName)

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

func (m *Mongo) ensureOrderIndexes(ctx context.Context) (err error) {
	col := m.DB.Collection(orderCollectionName)

	existingIndexes, err := m.existingIndexes(ctx, col)
	if err != nil {
		return err
	}

	indexesMap := map[string]mongo.IndexModel{
		"Pos Order Id": {
			Keys: bson.D{
				{Key: "pos_order_id", Value: 1},
			},
			Options: options.Index(),
		},
		"Order Time": {
			Keys: bson.D{
				{Key: "order_time.value", Value: -1},
			},
			Options: options.Index(),
		},
		"Unique Order": {
			Keys: bson.D{
				{Key: "order_id", Value: 1},
				{Key: "store_id", Value: 1},
				{Key: "delivery_service", Value: 1},
				{Key: "restaurant_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}

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
