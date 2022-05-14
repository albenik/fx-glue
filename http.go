package fxglue

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

func NewHTTPServerHook(srv *http.Server, log *zerolog.Logger) fx.Hook {
	return fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			go func() {
				if err = srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error().Err(err).Msg("Serve error!")
				}
			}()
			return nil
		},

		OnStop: srv.Shutdown,
	}
}
