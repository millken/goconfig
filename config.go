package goconfig

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"dario.cat/mergo"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-envconfig"
	"gopkg.in/yaml.v2"
)

// Option config option
type Option func(*option)

type option struct {
	path  string   //config file path
	files []string //config file names
	dirs  []string //config file dirs
}

// LoadWithDefault load config from file and env, if config file not exist, use default config
func LoadWithDefault[T any](t *T, opts ...Option) (*T, error) {
	var cfg = t
	var opt = new(option)
	for _, optFn := range opts {
		optFn(opt)
	}
	//load from default config file
	if opt.path == "" {
		for _, configDir := range opt.dirs {
			for _, configFile := range opt.files {
				dirPath, err := homedir.Expand(configDir)
				if err != nil {
					continue
				}
				path := filepath.Join(dirPath, configFile)
				if ok, _ := FileExists(path); ok {
					opt.path = path
					break
				}
			}
		}
	}
	if opt.path != "" {
		body, err := os.ReadFile(opt.path)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read config content")
		}
		if err = yaml.Unmarshal(body, cfg); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal config to struct")
		}
	}

	//load from env
	var envConfig = new(T)
	if err := envconfig.Process(context.Background(), envConfig); err != nil {
		return nil, errors.Wrap(err, "failed to process envconfig to struct")
	}
	if err := mergo.Merge(cfg, envConfig, mergo.WithOverride); err != nil {
		return nil, errors.Wrap(err, "failed to merge config")
	}
	return cfg, nil
}

// WithFile set config file path
func WithFile(path string) Option {
	return func(opt *option) {
		opt.path = path
	}
}

// Load config from file and env
func Load[T any](opts ...Option) (*T, error) {
	var cfg = new(T)
	return LoadWithDefault(cfg, opts...)
}

// File names from which we attempt to read configuration.
func SetDefaultConfigFiles(files ...string) Option {
	return func(opt *option) {
		opt.files = files
	}
}

// Launchd doesn't set root env variables, so there is default
func SetDefaultConfigDirs(dirs ...string) Option {
	return func(opt *option) {
		opt.dirs = dirs
	}
}

// GetCurrentDirectory return current directory
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// FileExists checks to see if a file exist at the provided path.
func FileExists(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			// ignore missing files
			return false, nil
		}
		return false, err
	}
	f.Close()
	return true, nil
}
