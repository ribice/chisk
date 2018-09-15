package zerolog

import (
	"context"
	"os"

	"github.com/ribice/chisk/model"

	"github.com/rs/zerolog"
)

// ZLog represents zerolog logger
type ZLog struct {
	logger *zerolog.Logger
}

// New instantiates new zero logger
func New() *ZLog {
	z := zerolog.New(os.Stdout)
	return &ZLog{
		logger: &z,
	}
}

// Log logs using zerolog
func (z *ZLog) Log(ctx context.Context, source, msg string, err error, params map[string]interface{}) {

	if params == nil {
		params = make(map[string]interface{})
	}

	params["source"] = source

	if user, ok := ctx.Value(chisk.KeyString("_authuser")).(*chisk.AuthUser); ok {
		params["id"] = user.ID
		params["username"] = user.DisplayName
	}

	if err != nil {
		params["error"] = err
		z.logger.Error().Fields(params).Msg(msg)
		return
	}

	z.logger.Info().Fields(params).Msg(msg)
}
