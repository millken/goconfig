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
	F float64  `yaml:"f" env:"F"`
}

func TestConfig(t *testing.T) {
	require := require.New(t)
	os.Setenv("B", "gg")
	os.Setenv("A", "1")
	defer func() {
		os.Unsetenv("B")
		os.Unsetenv("A")
	}()
	c1, err := Load[Config1]()
	require.NoError(err)
	require.Equal(c1.B, "gg")
	require.Equal(c1.A, 1)
	// require.Equal(c2.B, "gg")
}

func TestConfigWithDefault(t *testing.T) {
	require := require.New(t)
	c1 := Config1{
		A: 1,
		B: "a",
	}
	c2, err := LoadWithDefault[Config1](&c1)
	require.NoError(err)
	require.Equal(c2.A, 1)
	require.Equal(c2.B, "a")
}

func TestConfigWithFile(t *testing.T) {
	require := require.New(t)
	c2, err := Load[Config2](WithFile("testdata/test.yml"))
	require.NoError(err)
	require.Equal(123, c2.A)
	require.Equal("xxx", c2.B)
	require.Equal(true, c2.C)
	require.Equal([]string{"abc", "efg", "hij"}, c2.D)
	require.Equal(0.0, c2.F)
}

func TestConfigWithDefaultAndEnv(t *testing.T) {
	require := require.New(t)
	c1 := Config1{
		A: 1,
		B: "a",
	}
	os.Setenv("B", "gg")
	defer os.Unsetenv("B")
	c2, err := LoadWithDefault[Config1](&c1)
	require.NoError(err)
	require.Equal(1, c2.A)
	require.Equal("gg", c2.B)
}

func TestConfigWithDefaultAndFile(t *testing.T) {
	require := require.New(t)
	c1 := Config2{
		A: 1,
		F: 0.03,
	}
	c2, err := LoadWithDefault[Config2](&c1, WithFile("testdata/test.yml"))
	require.NoError(err)
	require.Equal(c2.A, 123)
	require.Equal(c2.B, "xxx")
	require.True(c2.C)
	require.Equal([]string{"abc", "efg", "hij"}, c2.D)
	require.Equal(0.03, c2.F)
}

func TestConfigWithDefaultAndEnvAndFile(t *testing.T) {
	require := require.New(t)
	c1 := Config2{
		A: 1,
		B: "a",
	}
	os.Setenv("B", "gg")
	os.Setenv("F", "0.111")
	defer func() {
		os.Unsetenv("B")
		os.Unsetenv("F")
	}()
	c2, err := LoadWithDefault[Config2](&c1, WithFile("testdata/test.yml"))
	require.NoError(err)
	require.Equal(123, c2.A)
	require.Equal("gg", c2.B)
	require.True(c2.C)
	require.Equal([]string{"abc", "efg", "hij"}, c2.D)
	require.Equal(0.111, c2.F)
}
