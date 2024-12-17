package selector

import "time"

type Delivery3plOrder struct {
	Id                 string
	DeliveryService    string
	DeliveryExternalID string
	CreatedTimeFrom    time.Time
	CreatedTimeTo      time.Time
	UpdatedTimeFrom    time.Time
	UpdatedTimeTo      time.Time
	Status             string
	VersionID          int
	CancelState        string
}

func EmptyDelivery3plSearch() Delivery3plOrder {
	return Delivery3plOrder{}
}

func (d Delivery3plOrder) SetStatus(status string) Delivery3plOrder {
	d.Status = status
	return d
}

func (d Delivery3plOrder) SetUpdatedTimeFrom(updatedTimeFrom time.Time) Delivery3plOrder {
	d.UpdatedTimeFrom = updatedTimeFrom
	return d
}

func (d Delivery3plOrder) SetUpdatedTimeTo(updatedTimeTo time.Time) Delivery3plOrder {
	d.UpdatedTimeTo = updatedTimeTo
	return d
}

func (d Delivery3plOrder) SetCreatedTimeFrom(createdTimeFrom time.Time) Delivery3plOrder {
	d.CreatedTimeFrom = createdTimeFrom
	return d
}

func (d Delivery3plOrder) SetCreatedTimeTo(createdTimeTo time.Time) Delivery3plOrder {
	d.CreatedTimeTo = createdTimeTo
	return d
}

func (d Delivery3plOrder) HasStatus() bool {
	return d.Status != ""
}

func (d Delivery3plOrder) HasUpdatedTimeFrom() bool {
	return d.UpdatedTimeFrom != time.Time{}
}

func (d Delivery3plOrder) HasUpdatedTimeTo() bool {
	return d.UpdatedTimeTo != time.Time{}
}

func (d Delivery3plOrder) HasCreatedTimeFrom() bool {
	return d.CreatedTimeFrom != time.Time{}
}

func (d Delivery3plOrder) HasCreatedTimeTo() bool {
	return d.CreatedTimeTo != time.Time{}
}
