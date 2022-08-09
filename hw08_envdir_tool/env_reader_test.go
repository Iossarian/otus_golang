package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReadDir(t *testing.T) {
	t.Run("test dir was read", func(t *testing.T) {
		env, err := ReadDir("testdata/env")

		expected := Environment{
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: true},
			"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
		}

		require.Equal(t, expected, env)
		require.Nil(t, err)
	})

	t.Run("test dir read failed", func(t *testing.T) {
		env, err := ReadDir("")

		require.Len(t, env, 0)
		require.NotNil(t, err)
	})
}
