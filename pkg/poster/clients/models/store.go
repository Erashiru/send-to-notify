package models

type GetSpotsResponse struct {
	Response []GetSpotsBody `json:"response"`
	ErrorResponse
}

type GetSpotsBody struct {
	SpotId     string            `json:"spot_id"`
	SpotName   string            `json:"spot_name"`
	SpotAdress string            `json:"spot_adress"`
	Storages   []GetSpotsStorage `json:"storages"`
}

type GetSpotsStorage struct {
	StorageId     int    `json:"storage_id"`
	StorageName   string `json:"storage_name"`
	StorageAdress string `json:"storage_adress"`
}
