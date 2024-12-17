package mongodb

import (
	"context"
	drivers2 "github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	connectionTimeout = 5 * time.Second
	ensureIdxTimeout  = 5 * time.Second

	storeCollectionName      = "restaurants"
	storeGroupCollectionName = "restaurant_groups"
	userStoreCollectionName  = "user_restaurants"
	storeTypeCollectionName  = "restaurant_types"
	apiTokenCollection       = "api_tokens"
	virtualStoreCollection   = "virtual_stores"
	tapRestaurantCollection  = "tap_restaurant"
)

type Mongo struct {
	connURL string
	dbName  string

	client *mongo.Client
	db     *mongo.Database

	connectionTimeout time.Duration
	ensureIdxTimeout  time.Duration

	storeRepo               drivers2.StoreRepository
	storeGroupRepo          drivers2.StoreGroupRepository
	userStoreRepo           drivers2.UserStoreRepository
	storeTypeRepo           drivers2.StoreTypeRepository
	apiTokenRepo            drivers2.ApiTokensRepository
	virtualStoreRepo        drivers2.VirtualRepository
	tapRestaurantRepository drivers2.TapRestaurantRepository
}

func New(conf drivers2.DataStoreConfig) (*Mongo, error) {
	if conf.URL == "" {
		return nil, errors.Wrap(drivers2.ErrInvalidConfigStruct, "conf.URL is nil")
	}

	if conf.DataBaseName == "" {
		return nil, errors.Wrap(drivers2.ErrInvalidConfigStruct, "conf.DataBaseName is nil")
	}

	return &Mongo{
		connURL:           conf.URL,
		dbName:            conf.DataBaseName,
		connectionTimeout: connectionTimeout,
		ensureIdxTimeout:  ensureIdxTimeout,
	}, nil
}

func (m *Mongo) Name() string {
	return "Mongo"
}

func (m *Mongo) Connect(cli *mongo.Client) error {

	var err error

	ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	defer cancel()

	if err = m.initializeClient(ctx, cli); err != nil {
		return err
	}

	if err = m.Ping(); err != nil {
		return err
	}
	m.db = m.client.Database(m.dbName)

	return m.ensureIndexes(ctx)
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

func (m *Mongo) StoreRepository() drivers2.StoreRepository {
	if m.storeRepo == nil {
		m.storeRepo = NewStoreRepository(m.db.Collection(storeCollectionName))
	}
	return m.storeRepo
}
func (m *Mongo) ApiTokensRepository() drivers2.ApiTokensRepository {
	if m.apiTokenRepo == nil {
		m.apiTokenRepo = NewApiTokensRepository(m.db.Collection(apiTokenCollection))
	}
	return m.apiTokenRepo
}

func (m *Mongo) StoreGroupRepository() drivers2.StoreGroupRepository {
	if m.storeGroupRepo == nil {
		m.storeGroupRepo = NewStoreGroupRepository(m.db.Collection(storeGroupCollectionName))
	}
	return m.storeGroupRepo
}

func (m *Mongo) UserStoreRepository() drivers2.UserStoreRepository {
	if m.userStoreRepo == nil {
		m.userStoreRepo = NewUserStoreRepository(m.db.Collection(userStoreCollectionName))
	}
	return m.userStoreRepo
}
func (m *Mongo) StoreTypeRepository() drivers2.StoreTypeRepository {
	if m.storeTypeRepo == nil {
		m.storeTypeRepo = NewStoreTypeRepository(m.db.Collection(storeTypeCollectionName))
	}
	return m.storeTypeRepo
}

func (m *Mongo) VirtualRepository() drivers2.VirtualRepository {
	if m.virtualStoreRepo == nil {
		m.virtualStoreRepo = NewVirtualRepository(m.db.Collection(virtualStoreCollection))
	}
	return m.virtualStoreRepo
}

func (m *Mongo) TapRestaurantRepository() drivers2.TapRestaurantRepository {
	if m.tapRestaurantRepository == nil {
		m.tapRestaurantRepository = NewTapRestaurantRepository(m.db.Collection(tapRestaurantCollection))
	}
	return m.tapRestaurantRepository
}

func (m *Mongo) ensureIndexes(ctx context.Context) error {

	col := m.db.Collection(storeCollectionName)

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
	defer closeCur(cur)

	res := make(map[string]struct{})

	for cur.Next(ctx) {
		res[cur.Current.Lookup("name").StringValue()] = struct{}{}
	}

	if err = cur.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func closeCur(cur *mongo.Cursor) {
	if err := cur.Close(context.Background()); err != nil {
		log.Err(err).Msg("closing cursor:")
	}
}

func toObjectIDs(ids []string) ([]primitive.ObjectID, error) {
	var objIDs []primitive.ObjectID
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objIDs = append(objIDs, oid)
	}
	return objIDs, nil
}
