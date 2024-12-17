package models

import "time"

type Newsletter struct {
	Name              string    `bson:"name"`
	RestaurantGroupId string    `bson:"restaurant_group_id"`
	Text              string    `bson:"text"`
	Recipients        []string  `bson:"recipients"`
	Status            string    `bson:"status"`
	SendDate          time.Time `bson:"send_date"`
	CreatedAt         time.Time `bson:"created_at"`
	UpdatedAt         time.Time `bson:"updated_at"`
}
