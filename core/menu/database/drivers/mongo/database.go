package mongo

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models/custom"
	"github.com/kwaaka-team/orders-core/service/entity_changes_history"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

const (
	connectionTimeout = 15 * time.Second
	ensureIdxTimeout  = 300 * time.Second

	menuCollection            = "menus"
	storeCollection           = "restaurants"
	menuUploadTransaction     = "menu_upload_transactions"
	ordersCollection          = "orders"
	userRestaurantsCollection = "user_restaurants"
	stoplistTransaction       = "stoplist_transaction"
	sequencesCollection       = "sequences"
	promoCollection           = "promo"
	msPositionsCollection     = "moysklad_positions"
	bkOffers                  = "bk_offers"
	restGroupMenu             = "restaurant_group_menu"
)

type Mongo struct {
	connURL string
	dbName  string

	client *mongo.Client
	DB     *mongo.Database

	connectionTimeout time.Duration
	ensureIdxTimeout  time.Duration

	menuRepo        drivers.MenuRepository
	storeRepo       drivers.StoreRepository
	menuUTRepo      drivers.MenuUploadTransactionRepository
	stopListRepo    drivers.StopListTransactionRepository
	sequencesRepo   drivers.SequencesRepository
	promoRepo       drivers.PromoRepository
	msPositionsRepo drivers.MSPositionsRepository
	bkOffersRepo    drivers.BkOffersRepository
	restGroupMenu   drivers.RestaurantGroupMenuRepository
}

func (m *Mongo) Name() string { return "Mongo" }

func (m *Mongo) Client() *mongo.Client {
	return m.client
}

func (m *Mongo) Connect(client *mongo.Client) error {

	var err error

	ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	defer cancel()

	if client != nil {
		m.client = client
	}

	if client == nil {

		m.client, err = mongo.Connect(ctx, options.Client().ApplyURI(m.connURL))
		if err != nil {
			return err
		}

	}

	if err = m.Ping(); err != nil {
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

func (m *Mongo) StoreRepository() drivers.StoreRepository {
	if m.storeRepo == nil {
		m.storeRepo = NewStoreRepository(
			m.DB.Collection(storeCollection),
		)
	}
	return m.storeRepo
}

func (m *Mongo) MenuRepository(entityChangesHistoryRepo entity_changes_history.Repository) drivers.MenuRepository {
	if m.menuRepo == nil {
		m.menuRepo = NewMenuRepository(
			m.DB.Collection(menuCollection), entityChangesHistoryRepo,
		)
	}
	return m.menuRepo
}

func (m *Mongo) RestGroupMenuRepository() drivers.RestaurantGroupMenuRepository {
	if m.restGroupMenu == nil {
		m.restGroupMenu = NewRestaurantGroupMenuRepository(
			m.DB.Collection(restGroupMenu),
		)
	}

	return m.restGroupMenu
}

func (m *Mongo) MenuUploadTransactionRepository() drivers.MenuUploadTransactionRepository {
	if m.menuUTRepo == nil {
		m.menuUTRepo = NewMenuUploadTransaction(
			m.DB.Collection(menuUploadTransaction),
		)
	}
	return m.menuUTRepo
}

func (m *Mongo) BkOffersRepository() drivers.BkOffersRepository {
	if m.bkOffersRepo == nil {
		m.bkOffersRepo = NewBkOffers(
			m.DB.Collection(bkOffers),
		)
	}
	return m.bkOffersRepo
}

func (m *Mongo) MSPositionsRepository() drivers.MSPositionsRepository {
	if m.msPositionsRepo == nil {
		m.msPositionsRepo = NewMSPositionsRepository(
			m.DB.Collection(msPositionsCollection),
		)
	}
	return m.msPositionsRepo
}

func (m *Mongo) StopListTransactionRepository() drivers.StopListTransactionRepository {
	if m.stopListRepo == nil {
		m.stopListRepo = NewStopListTransaction(
			m.DB.Collection(stoplistTransaction),
		)
	}
	return m.stopListRepo
}

func (m *Mongo) SequencesRepository() drivers.SequencesRepository {
	if m.sequencesRepo == nil {
		m.sequencesRepo = NewSequencesRepository(m.DB.Collection(sequencesCollection))
	}
	return m.sequencesRepo
}

func (m *Mongo) PromoRepository() drivers.PromoRepository {
	if m.promoRepo == nil {
		m.promoRepo = NewPromoRepository(m.DB.Collection(promoCollection))
	}
	return m.promoRepo
}

func (m *Mongo) DataBase() *mongo.Database {
	return m.DB
}

func New(conf drivers.DataStoreConfig) (drivers.DataStore, error) {
	if conf.URL == "" {
		return nil, drivers.ErrInvalidConfigStruct
	}

	if conf.DataBaseName == "" {
		return nil, drivers.ErrInvalidConfigStruct
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
	if err := m.ensureMenuUploadTransactionIndexes(ctx); err != nil {
		return err
	}

	return nil
}

func (m *Mongo) ensureMenuIndexes(ctx context.Context) (err error) {
	col := m.DB.Collection(menuCollection)

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

func (m *Mongo) ensureMenuUploadTransactionIndexes(ctx context.Context) (err error) {
	col := m.DB.Collection(menuUploadTransaction)

	existingIndexes, err := m.existingIndexes(ctx, col)
	if err != nil {
		return err
	}

	indexesMap := map[string]mongo.IndexModel{
		"store_id_idx": {
			Keys:    bson.D{{Key: "restaurant_id", Value: 1}},
			Options: options.Index(),
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

func (m *Mongo) StartSession(ctx context.Context) (context.Context, drivers.TxCallback, error) {
	wc := writeconcern.Majority()
	rc := readconcern.Snapshot()
	txOpts := options.Transaction().
		SetWriteConcern(wc).
		SetReadConcern(rc).
		SetReadPreference(readpref.Primary())

	session, err := m.client.StartSession()
	if err != nil {
		return nil, nil, err
	}

	if err = session.StartTransaction(txOpts); err != nil {
		return nil, nil, err
	}

	return mongo.NewSessionContext(ctx, session), callback(session), nil
}

func callback(session mongo.Session) func(err error) error {
	return func(err error) error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		defer session.EndSession(ctx)

		if err == nil {
			err = session.CommitTransaction(ctx)
		}
		if err != nil {
			if abortErr := session.AbortTransaction(ctx); abortErr != nil {
				var errs custom.Error
				errs.Append(err, abortErr)
				err = errs.ErrorOrNil()
			}
		}
		return err
	}
}
