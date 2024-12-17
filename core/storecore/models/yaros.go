package models

type StoreYarosConfig struct {
	Username   string `bson:"username" json:"username"`
	Password   string `bson:"password" json:"password"`
	StoreId    string `bson:"store_id" json:"store_id"`
	Department string `bson:"department" json:"department"`
}
