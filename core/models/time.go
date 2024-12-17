package models

import (
	"encoding/json"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type Time struct {
	time.Time
}

func TimeNow() Time {
	return Time{
		Time: time.Now(),
	}
}

func FromTime(t time.Time) Time {
	return Time{
		Time: t,
	}
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	err = json.Unmarshal(data, &t.Time)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.Round(time.Second))
}

func (t *Time) UnmarshalBSONValue(typ bsontype.Type, data []byte) (err error) {
	rv := bson.RawValue{Type: typ, Value: data}
	err = rv.Unmarshal(&t.Time)
	t.Time = t.Time.UTC()
	return
}

func (t Time) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(t.Time.UTC())
}

type TransactionTime struct {
	RestaurantID string `bson:"restaurant_id" json:"restaurant_id"`
	Value        Time   `bson:"value" json:"value"`
	TimeZone     string `bson:"tz" json:"tz"`
	TimeStamp    int64  `bson:"-" json:"timestamp"`
	UTCOffset    string `bson:"utc_offset" json:"utc_offset"`
}

func (tt *TransactionTime) GetLocalTime() (time.Time, error) {
	offsetMinutes, err := strconv.Atoi(tt.TimeZone)
	if err == nil {
		offset := time.Duration(offsetMinutes) * time.Minute
		return tt.Value.Add(offset), nil
	}

	loc, err := time.LoadLocation(tt.TimeZone)
	if err != nil {
		return time.Time{}, err
	}

	return tt.Value.In(loc), nil
}
