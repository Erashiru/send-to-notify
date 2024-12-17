package dto

import "time"

type StoreManageRequest struct {
	StoreID string    `json:"store_id"`
	Until   time.Time `json:"until"`
}

type StoreStatusResponse struct {
	StoreID  string `json:"store_id"`
	IsActive bool   `json:"is_active"`
}

type StoreScheduleResponse struct {
	Timezone string          `json:"timezone,omitempty" bson:"timezone,omitempty"`
	Schedule []StoreSchedule `json:"schedule,omitempty" bson:"schedule,omitempty"`
}

type StoreSchedule struct {
	DayOfWeek int                `json:"day_of_week" bson:"day_of_week"`
	TimeSlots []ScheduleTimeSlot `json:"time_slots" bson:"time_slots"`
}

type ScheduleTimeSlot struct {
	Opening string `json:"opening" bson:"opening"`
	Closing string `json:"closing" bson:"closing"`
}
