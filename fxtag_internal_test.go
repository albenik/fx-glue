package fxglue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTag(t *testing.T) {
	t.Parallel()

	t.Run("Supply", func(t *testing.T) {
		t.Parallel()

		anno, name, err := parseTag("supply")
		if assert.NoError(t, err) {
			assert.Empty(t, anno)
			assert.Empty(t, name)
		}
	})

	t.Run("SupplyNameFoo", func(t *testing.T) {
		t.Parallel()

		anno, name, err := parseTag("supply,name=foo")
		if assert.NoError(t, err) {
			assert.Equal(t, "name", anno)
			assert.Equal(t, "foo", name)
		}
	})

	t.Run("SupplyGroupBar", func(t *testing.T) {
		t.Parallel()

		anno, name, err := parseTag("supply,group=bar")
		if assert.NoError(t, err) {
			assert.Equal(t, "group", anno)
			assert.Equal(t, "bar", name)
		}
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		testCases := []struct {
			Name  string
			Value string
		}{
			{Name: "Empty", Value: ""},
			{Name: "Invalid", Value: "invalid"},
			{Name: "InvalidBroken1", Value: "invalid,"},
			{Name: "InvalidBroken2", Value: "invalid,name"},
			{Name: "InvalidBroken3", Value: "invalid,name-foo"},
			{Name: "InvalidAnnotation", Value: "supply,foo=bar"},
			{Name: "InvalidNameFoo", Value: "invalid,name=foo"},
			{Name: "InvalidGroupBar", Value: "invalid,group=bar"},
		}

		for _, c := range testCases {
			c := c
			t.Run(c.Name, func(t *testing.T) {
				t.Parallel()

				_, _, err := parseTag(c.Value)
				assert.Error(t, err)
			})
		}
	})
}
