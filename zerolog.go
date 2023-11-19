package fxglue

import (
	"fmt"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

// fxLogger is a App event logger that logs events to Zerolog.
// Adaped from `go.uber.org/fx/fxevent.ZapLogger`.
type fxLogger struct {
	logger *zerolog.Logger
}

func WithZerolog(log *zerolog.Logger) fx.Option {
	return fx.Options(
		fx.WithLogger(func() fxevent.Logger { return &fxLogger{logger: log} }),
		fx.Supply(log),
	)
}

// LogEvent logs the given event to the provided Zap logger.
func (l *fxLogger) LogEvent(event fxevent.Event) { //nolint:funlen,cyclop
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.logger.Info().Str("callee", e.FunctionName).Str("caller", e.CallerName).Msg("OnStart executing")

	case *fxevent.OnStartExecuted:
		log := l.logger.With().Str("callee", e.FunctionName).Str("caller", e.CallerName).
			Stringer("runtime", e.Runtime).Logger()
		if e.Err != nil {
			log.Error().Err(e.Err).Msg("OnStart error!")
		} else {
			log.Info().Msg("OnStart executed")
		}

	case *fxevent.OnStopExecuting:
		l.logger.Info().Str("callee", e.FunctionName).Str("caller", e.CallerName).Msg("OnStop executing")

	case *fxevent.OnStopExecuted:
		log := l.logger.With().Str("callee", e.FunctionName).Str("caller", e.CallerName).
			Stringer("runtime", e.Runtime).Logger()
		if e.Err != nil {
			log.Error().Err(e.Err).Msg("OnStop error!")
		} else {
			l.logger.Info().Msg("OnStop executing")
		}

	case *fxevent.Started:
		if e.Err != nil {
			l.logger.Error().Msg("Start failed!")
		} else {
			l.logger.Info().Msg("Started")
		}

	case *fxevent.Stopping:
		l.logger.Info().Stringer("signal", e.Signal).Msg("OS signal received")

	case *fxevent.Stopped:
		if e.Err != nil {
			l.logger.Error().Err(e.Err).Msg("Stopped with errors!")
		} else {
			l.logger.Info().Msg("Stopped")
		}

	case *fxevent.Supplied:
		log := l.logger.With().Str("type", e.TypeName).Logger()
		if e.Err != nil {
			log.Error().Err(e.Err).Msg("Supply error!")
		} else {
			log.Info().Msg("Supplied")
		}

	case *fxevent.Provided:
		log := l.logger.With().Str("contructor", e.ConstructorName).Strs("types", e.OutputTypeNames).Logger()
		if e.Err != nil {
			log.Error().Err(e.Err).Msg("Provide error!")
		} else {
			log.Info().Msg("Provided")
		}

	case *fxevent.Invoking:
		// Do not log stack as it will make logs hard to read.
		l.logger.Info().Str("function", e.FunctionName).Msg("Invoking")

	case *fxevent.Invoked:
		log := l.logger.With().Str("function", e.FunctionName).Str("stack", e.Trace).Logger()
		if e.Err != nil {
			log.Error().Err(e.Err).Msg("Invocation error!")
		} else {
			log.Info().Msg("Invoked")
		}

	case *fxevent.RollingBack:
		l.logger.Error().Err(e.StartErr).Msg("Start failed, performing rollback")

	case *fxevent.RolledBack:
		if e.Err != nil {
			l.logger.Error().Err(e.Err).Msg("Rollback failed!")
		} else {
			l.logger.Info().Msg("Rollback performed")
		}

	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.logger.Error().Err(e.Err).Msg("Custom logger error!")
		} else {
			l.logger.Info().Str("function", e.ConstructorName).Msg("Custom logger")
		}

	default:
		l.logger.Error().Str("fxevent", fmt.Sprintf("%T", event)).Msg("Unknown fxevent type")
	}
}
