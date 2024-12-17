package models

import (
	coreOrderModels "github.com/kwaaka-team/orders-core/core/models"
)

func (req Event) ToOrderModel() (coreOrderModels.Order, error) {
	//total, err := strconv.ParseFloat(req.Data.Amount, 64)
	//if err != nil {
	//	return coreOrderModels.Order{}, fmt.Errorf("parse req.Data.Amount string to float64 error: %s", err.Error())
	//}

	order := coreOrderModels.Order{
		RestaurantID: req.RestaurantId,
		PosOrderID:   req.Data.OrderId,
		PosType:      "jowi",
		Courier: coreOrderModels.Courier{
			Name:        req.Data.CourierName,
			PhoneNumber: req.Data.CourierPhone,
		},
		//TotalCustomerToPay: coreOrderModels.Price{
		//	Value: total,
		//},
	}

	switch req.Status {
	case 0:
		order.Status = "NEW"
	case 1:
		order.Status = "ACCEPTED"
	case 2:
		order.Status = "CANCELLED"
	case 3:
		order.Status = "COOKING_STARTED"
	case 4:
		order.Status = "CLOSED"
	}

	return order, nil
}
