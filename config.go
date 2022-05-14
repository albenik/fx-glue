package fxglue

import (
	"errors"
	"reflect"
	"time"

	"go.uber.org/fx"
)

var ErrFxApplicationConfigMissing = errors.New("fx application config missing")

type AppConfig struct {
	StartTimeout time.Duration `env:"APP_START_TIMEOUT"`
	StopTimeout  time.Duration `env:"APP_STOP_TIMEOUT"`
}

func SupplyConfig(conf interface{}) fx.Option {
	items := make([]interface{}, 0)

	appConf := enumerateFields(&items, reflect.ValueOf(conf))
	if appConf == nil {
		return fx.Error(ErrFxApplicationConfigMissing)
	}

	return fx.Options(
		fx.Options(
			fx.StartTimeout(appConf.StartTimeout),
			fx.StopTimeout(appConf.StopTimeout),
		),
		fx.Supply(items...),
	)
}

func enumerateFields(list *[]interface{}, v reflect.Value) *AppConfig {
	var conf *AppConfig

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return conf
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldVal := field.Interface()

		switch c := fieldVal.(type) {
		case AppConfig:
			conf = &c
			continue
		case *AppConfig:
			conf = c
			continue
		}

		anno, name, err := parseTag(t.Field(i).Tag.Get("fx"))
		if err != nil {
			panic(err)
		}

		if anno == "" {
			if c := enumerateFields(list, field); c != nil && conf == nil {
				conf = c
			}
			continue
		}

		switch anno {
		case "name":
			*list = append(*list, fx.Annotated{Name: name, Target: fieldVal})
		case "group":
			*list = append(*list, fx.Annotated{Group: name, Target: fieldVal})
		default:
			*list = append(*list, fieldVal)
		}
	}

	return conf
}
