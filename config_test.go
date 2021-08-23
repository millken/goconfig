package goconfig

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type Config1 struct {
	A int    `yaml:"a" env:"A"`
	B string `yaml:"b" env:"B"`
}
type Config2 struct {
	A int      `yaml:"a" env:"A"`
	B string   `yaml:"b" env:"B"`
	C bool     `yaml:"c"`
	D []string `yaml:"d"`
}

func TestLoadConfigFromEnv(t *testing.T) {
	require := require.New(t)
	c1 := Config1{
		A: 1,
		B: "a",
	}
	var c2 Config1
	os.Setenv("B", "gg")
	err := LoadConfigFromEnv(&c1, &c2)
	require.NoError(err)
	require.Equal(c1.B, "gg")
	require.Equal(c2.B, "gg")
}

func TestLoadConfigFromFile(t *testing.T) {
	require := require.New(t)
	c1 := Config2{
		A: 1,
		B: "a",
	}
	var c2 Config2
	err := LoadConfigFromFile("testdata/test.yml", &c1, &c2)
	require.NoError(err)
	require.Equal(c1.B, "f3a")
	require.False(c2.C)
	require.True(c1.C)
}
