package rx

import (
	"github.com/rs/zerolog/log"
)

func GetErrFuncLogError(header string) func(error) interface{} {
	return func(err error) interface{} {
		log.Error().Err(err).Str("header", header).Msg("")
		return err
	}
}
