package models

type StoreType struct {
	ID     string   `bson:"_id" json:"id"`
	Type   string   `bson:"type" json:"type"`
	Format []string `bson:"format" json:"format"`
}
