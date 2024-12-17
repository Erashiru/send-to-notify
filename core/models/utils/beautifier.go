package utils

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

func Beautify(message string, model any) {
	body, err := json.Marshal(model)
	if err != nil {
		log.Err(err).Msgf("[beautify] %s marshal error %v", message, model)
		return
	}

	log.Info().Msgf("[beautify] %s marshal body: %s", message, string(body))
}

func PointerOfFloat(val float64) *float64 {
	return &val
}

func GetJsonFormatFromModel(req interface{}) (string, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
