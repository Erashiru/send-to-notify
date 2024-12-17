package utils

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
)

func Beautify(message string, model any) {
	body, err := json.Marshal(model)
	if err != nil {
		log.Err(err).Msgf("[beautify] %s marshal error", message)
		return
	}

	log.Info().Msgf("[beautify] %s marshal body: %s", message, string(body))
}
