package rx

import (
	"context"
	"github.com/rs/zerolog/log"
)

//GetSideEffectLog simply log the input with a header string
func GetSideEffectLog(header string) func(context.Context, interface{}) (interface{}, error) {
	return func(_ context.Context, i interface{}) (interface{}, error) {
		if len(header) == 0 {
			header = "log"
		}
		log.Info().Interface(header, i).Msg("")
		return i, nil
	}
}
