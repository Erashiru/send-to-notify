package models

type SetScheduleRequest struct {
	Timezone  string     `json:"timezone"`
	Schedules []Schedule `json:"schedule"`
}

type Schedule struct {
	DayOfWeek int        `json:"day_of_week"`
	TimeSlots []TimeSlot `json:"time_slots"`
}

type TimeSlot struct {
	Opening string `json:"opening"`
	Closing string `json:"closing"`
}
