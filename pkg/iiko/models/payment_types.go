package models

import "github.com/kwaaka-team/orders-core/core/storecore/models"

type PaymentTypes struct {
	CASH    models.PaymentType
	DELAYED models.PaymentType
}
