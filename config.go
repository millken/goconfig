package goconfig

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/imdario/mergo"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-envconfig"
	"gopkg.in/yaml.v2"
)

//LoadConfigFromEnv loaded config from environment variables, override default config
func LoadConfigFromEnv(defaultConfig, envConfig interface{}) error {

	if err := envconfig.Process(context.Background(), envConfig); err != nil {
		return errors.Wrap(err, "failed to process envconfig to struct")
	}

	if err := mergo.Merge(defaultConfig, envConfig, mergo.WithOverride); err != nil {
		return errors.Wrap(err, "failed to merge config")
	}
	return nil
}

func LoadConfigFromFile(path string, defaultConfig, envConfig interface{}) error {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "failed to read config content")
	}
	if err = yaml.Unmarshal(body, defaultConfig); err != nil {
		return errors.Wrap(err, "failed to unmarshal config to struct")
	}

	return LoadConfigFromEnv(defaultConfig, envConfig)
}

var (
	defaultConfigFiles []string
	defaultConfigDirs  []string
)

// File names from which we attempt to read configuration.
func SetDefaultConfigFiles(files ...string) {
	defaultConfigFiles = files
}

// Launchd doesn't set root env variables, so there is default
func SetDefaultConfigDirs(dirs ...string) {
	defaultConfigDirs = dirs
}

//GetCurrentDirectory return current directory
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
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

// FindDefaultConfigPath returns the first path that contains a config file.
// If none of the combination of DefaultConfigDirs and DefaultConfigFiles
// contains a config file, return empty string.
func FindDefaultConfigPath() string {
	for _, configDir := range defaultConfigDirs {
		for _, configFile := range defaultConfigFiles {
			dirPath, err := homedir.Expand(configDir)
			if err != nil {
				continue
			}
			path := filepath.Join(dirPath, configFile)
			if ok, _ := FileExists(path); ok {
				return path
			}
		}
	}
	return ""
}
