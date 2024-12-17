package repository

import (
	"errors"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	"github.com/kwaaka-team/orders-core/service/legalentity/models"
	leSel "github.com/kwaaka-team/orders-core/service/legalentity/models/selector"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func errorSwitch(err error) error {
	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		return models.ErrNotFound
	case mongo.IsDuplicateKeyError(err):
		return models.ErrDuplicateData
	default:
		return err
	}
}

func (r *Repo) getObjectIDFromString(objID string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(objID)
}

func (r *Repo) filterFrom(query leSel.LegalEntityForm) bson.D {
	var res bson.D

	if query.HasName() {
		res = append(res, bson.E{Key: "name", Value: query.Name})
	}

	if query.HasBIN() {
		res = append(res, bson.E{Key: "bin", Value: query.BIN})
	}

	if query.HasKNP() {
		res = append(res, bson.E{Key: "knp", Value: query.KNP})
	}

	if query.HasPaymentType() {
		res = append(res, bson.E{Key: "payment_type", Value: query.PaymentType})
	}

	if query.HasLinkedAccManager() {
		res = append(res, bson.E{Key: "linked_acc_manager", Value: query.LinkedAccManager})
	}

	if query.HasSalesID() {
		res = append(res, bson.E{Key: "sales_id", Value: query.SalesID})
	}

	if query.HasContacts() {
		res = append(res, bson.E{Key: "contacts", Value: query.Contacts})
	}

	if query.HasSalesComment() {
		res = append(res, bson.E{Key: "sales_comment", Value: query.SalesComment})
	}

	if query.HasStoreIds() {
		res = append(res, bson.E{Key: "store_ids", Value: query.StoreIds})
	}

	if query.HasPaymentCycle() {
		res = append(res, bson.E{Key: "payment_cycle", Value: query.PaymentCycle})
	}

	return res
}

func applyPagination(pagination selector.Pagination) []bson.D {
	var withPagination []bson.D

	withPagination = append(withPagination, bson.D{{Key: "$skip", Value: pagination.Skip()}})
	withPagination = append(withPagination, bson.D{{Key: "$limit", Value: pagination.Limit}})

	return withPagination
}
