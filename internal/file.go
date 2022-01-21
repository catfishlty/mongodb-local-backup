package internal

import (
	"encoding/json"
	"errors"
	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

// Exist check if the given path is exists of file or directory
func Exist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// IsDir check if the given path is a directory or not
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// Mkdir create a new directory from the given path
func Mkdir(path string) error {
	return os.MkdirAll(path, 0777)
}

// IsFile check if the given path is a file or not
func IsFile(path string) bool {
	return !IsDir(path)
}

// FileTrans move file form `sourcePath` to `targetPath`
func FileTrans(sourcePath, targetPath string) error {
	b, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(targetPath, b, 0777)
	if err != nil {
		return err
	}
	err = os.Remove(sourcePath)
	if err != nil {
		return err
	}
	return nil
}

// ReadConfig read config file from given path and file format, support `json`, `yaml` and `toml`
func ReadConfig(file, format string) (*Config, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	switch format {
	case "json":
		return ReadJson(b)
	case "yaml":
		return ReadYaml(b)
	case "toml":
		return ReadToml(b)
	}
	return nil, errors.New("unsupported config file format")
}

func ReadJson(b []byte) (*Config, error) {
	var config Config
	err := json.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func ReadYaml(b []byte) (*Config, error) {
	var config Config
	err := yaml.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func ReadToml(b []byte) (*Config, error) {
	var config Config
	err := toml.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
