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


func Contains(arr []any, el any) bool {
	for _, v := range arr {
		if v == el {
			return true
		}
	}

	return false
}
