package repository

import (
	"context"
	errorsGo "errors"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	"github.com/kwaaka-team/orders-core/service/legalentity/models"
	leSel "github.com/kwaaka-team/orders-core/service/legalentity/models/selector"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	legalEntityCollection = "legal_entity"
	storesCollection      = "restaurants"
	brandsCollection      = "restaurant_groups"
	managersCollection    = "managers"
	salesCollection       = "sales"

	ACTIVE   = "active"
	INACTIVE = "inactive"
	DISABLED = "disabled"
)

type Repo struct {
	legalEntityColl *mongo.Collection
}

func NewLegalEntityRepo(db *mongo.Database) *Repo {
	legalEntityColl := db.Collection(legalEntityCollection)

	return &Repo{legalEntityColl: legalEntityColl}
}

func (r *Repo) Insert(ctx context.Context, req models.LegalEntityForm) (string, error) {
	req.CreatedAt = time.Now().UTC()
	req.UpdatedAt = time.Now().UTC()
	req.Status = INACTIVE

	res, err := r.legalEntityColl.InsertOne(ctx, req)
	if err != nil {
		return "", errorSwitch(err)
	}

	objID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", models.ErrInvalidID
	}

	return objID.Hex(), nil
}

func (r *Repo) GetByID(ctx context.Context, id string) (models.LegalEntityView, error) {
	objectID, err := r.getObjectIDFromString(id)
	if err != nil {
		return models.LegalEntityView{}, models.ErrInvalidID
	}

	var legalEntity models.LegalEntityView
	pipeline := r.getLegalEntityView(objectID)

	cursor, err := r.legalEntityColl.Aggregate(ctx, pipeline)
	if err != nil {
		return models.LegalEntityView{}, errorSwitch(err)
	}
	defer errors.CloseCur(cursor)

	for cursor.Next(ctx) {
		if err := cursor.Decode(&legalEntity); err != nil {
			return models.LegalEntityView{}, errorSwitch(err)
		}
	}
	if err := cursor.Err(); err != nil {
		return models.LegalEntityView{}, errorSwitch(err)
	}

	if legalEntity.ID == "" {
		return models.LegalEntityView{}, models.ErrNotFound
	}

	return legalEntity, nil
}

func (r *Repo) Update(ctx context.Context, u leSel.LegalEntityForm, id string) error {
	objectID, err := r.getObjectIDFromString(id)
	if err != nil {
		return errorSwitch(err)
	}

	toUpdate := r.filterFrom(u)

	toUpdate = append(toUpdate, bson.E{Key: "updated_at", Value: time.Now().UTC()})

	filter := bson.M{"_id": objectID}

	update := bson.D{{Key: "$set", Value: toUpdate}}

	res, err := r.legalEntityColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return errorSwitch(err)
	}
	if res.ModifiedCount == 0 {
		return models.ErrInvalidID
	}

	return nil
}

func (r *Repo) Disable(ctx context.Context, id string) error {
	objectID, err := r.getObjectIDFromString(id)
	if err != nil {
		return models.ErrInvalidID
	}

	filter := bson.M{"_id": objectID}

	update := bson.M{
		"$set": bson.M{
			"status":     DISABLED,
			"updated_at": time.Now().UTC(),
		},
	}

	res, err := r.legalEntityColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return errorSwitch(err)
	}
	if res.ModifiedCount == 0 {
		return models.ErrInvalidID
	}

	return nil
}

func (r *Repo) List(ctx context.Context, pagination selector.Pagination, filter models.Filter) ([]models.GetListOfLegalEntitiesDB, error) {
	var legalEntities []models.GetListOfLegalEntitiesDB

	pipeline := r.getListOfLegalEntities(pagination, filter)

	cursor, err := r.legalEntityColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer errors.CloseCur(cursor)

	for cursor.Next(ctx) {
		var le models.GetListOfLegalEntitiesDB
		if err := cursor.Decode(&le); err != nil {
			return nil, err
		}
		legalEntities = append(legalEntities, le)
	}

	return legalEntities, nil
}

func (r *Repo) GetStores(ctx context.Context, pagination selector.Pagination, id string) (models.GetListOfStoresDB, error) {
	objectID, err := r.getObjectIDFromString(id)
	if err != nil {
		return models.GetListOfStoresDB{}, models.ErrInvalidID
	}

	var storesInfo models.GetListOfStoresDB

	pipeline := r.getListOfStores(pagination, objectID)

	cursor, err := r.legalEntityColl.Aggregate(ctx, pipeline)
	if err != nil {
		return models.GetListOfStoresDB{}, err
	}
	defer errors.CloseCur(cursor)

	for cursor.Next(ctx) {
		if err := cursor.Decode(&storesInfo); err != nil {
			return models.GetListOfStoresDB{}, err
		}
	}
	if cursor.Err() != nil {
		return models.GetListOfStoresDB{}, err
	}

	return storesInfo, nil
}

func (r *Repo) InsertDocument(ctx context.Context, id string, newDoc models.Document) error {
	objectID, err := r.getObjectIDFromString(id)
	if err != nil {
		return models.ErrInvalidID
	}

	newDoc.Status = ACTIVE

	filter := bson.M{"_id": objectID}

	update := bson.M{
		"$push": bson.M{
			"documents": newDoc,
		},
	}

	res, err := r.legalEntityColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return errorsGo.New("no legal entity found to update (InsertDocument)")
	}

	return nil
}

func (r *Repo) DisableDocument(ctx context.Context, legalEntityID, documentID string) error {
	legalEntityObjID, err := r.getObjectIDFromString(legalEntityID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": legalEntityObjID, "documents.id": documentID}
	update := bson.M{
		"$set": bson.M{"documents.$[elem].status": DISABLED},
	}
	arrayFilter := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"elem.id": documentID},
		},
	}
	updateOptions := options.Update().SetArrayFilters(arrayFilter)

	res, err := r.legalEntityColl.UpdateOne(ctx, filter, update, updateOptions)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errorsGo.New("no matching legal entity found")
	}

	if res.ModifiedCount == 0 {
		return errorsGo.New("no documents were updated")
	}

	return nil
}

func (r *Repo) GetAllDocumentsByLegalEntityID(ctx context.Context, pagination selector.Pagination, filter models.DocumentFilter, legalEntityID string) ([]models.Document, error) {
	legalEntityObjID, err := r.getObjectIDFromString(legalEntityID)
	if err != nil {
		return nil, err
	}

	pipeline := r.getAllDocumentsByLegalEntityID(legalEntityObjID, pagination, filter)

	cur, err := r.legalEntityColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer errors.CloseCur(cur)

	var documents []models.Document
	for cur.Next(ctx) {
		var doc models.UnwindDocumentsField
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		documents = append(documents, doc.Document)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	if len(documents) == 0 {
		return nil, errorsGo.New("no documents were found")
	}

	return documents, nil
}

func (r *Repo) GetDocumentDownloadLink(ctx context.Context, legalEntityID, documentID string) (string, error) {
	var documents models.GetDocumentByLegalEntityIDRequest

	legalEntityObjID, err := r.getObjectIDFromString(legalEntityID)
	if err != nil {
		return "", err
	}

	filter := bson.M{"_id": legalEntityObjID, "documents.id": documentID}
	projection := bson.M{"_id": 0, "documents.$": 1}
	opts := options.FindOne().SetProjection(projection)

	if err := r.legalEntityColl.FindOne(ctx, filter, opts).Decode(&documents); err != nil {
		if errorsGo.Is(err, mongo.ErrNoDocuments) {
			return "", errorSwitch(err)
		}
		return "", err
	}

	return documents.Documents[0].S3Link, nil
}
