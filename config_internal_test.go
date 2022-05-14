package fxglue

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

type Config struct {
	App     AppConfig
	Nested1 *NestedConfig
	Nested2 *NestedConfig `fx:"supply"`
}

type NestedConfig struct {
	Foo string `fx:"supply"`
	Bar string `fx:"supply,name=bar"`
	Baz string `fx:"supply,group=baz"`
}

func TestSupplyConfig(t *testing.T) {
	t.Parallel()

	conf := &Config{
		App: AppConfig{
			StartTimeout: 2 * time.Second,
			StopTimeout:  3 * time.Second,
		},
		Nested1: &NestedConfig{}, //nolint:exhaustruct
		Nested2: &NestedConfig{}, //nolint:exhaustruct
	}

	items := make([]interface{}, 0)
	appConf := enumerateFields(&items, reflect.ValueOf(conf))

	assert.Equal(t, &conf.App, appConf)
	assert.Equal(t, []interface{}{
		conf.Nested1.Foo,
		fx.Annotated{Target: conf.Nested1.Bar, Name: "bar"},
		fx.Annotated{Target: conf.Nested1.Baz, Group: "baz"},
		conf.Nested2,
	}, items)
}
