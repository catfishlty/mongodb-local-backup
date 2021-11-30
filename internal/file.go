package internal

import (
    "encoding/json"
    "errors"
    "github.com/BurntSushi/toml"
    "gopkg.in/yaml.v3"
    "io/ioutil"
    "os"
)

func Exist(path string) bool {
    _, err := os.Stat(path)
    if err != nil {
        if os.IsExist(err) {
            return true
        }
        return false
    }
    return true
}

func IsDir(path string) bool {
    s, err := os.Stat(path)
    if err != nil {
        return false
    }
    return s.IsDir()
}

func Mkdir(path string) error {
    return os.MkdirAll(path, 0777)
}

func IsFile(path string) bool {
    return !IsDir(path)
}

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

func ReadConfig(file, format string) (*Config, error) {
    b, err := ioutil.ReadFile(file)
    if err != nil {
        return nil, err
    }
    switch format {
    case "json":
        return readJson(b)
    case "yaml":
        return readYaml(b)
    case "toml":
        return readToml(b)
    }
    return nil, errors.New("unsupported config file format")
}

func readJson(b []byte) (*Config, error) {
    var config Config
    err := json.Unmarshal(b, &config)
    if err != nil {
        return nil, err
    }
    return &config, nil
}

func readYaml(b []byte) (*Config, error) {
    var config Config
    err := yaml.Unmarshal(b, &config)
    if err != nil {
        return nil, err
    }
    return &config, nil
}

func readToml(b []byte) (*Config, error) {
    var config Config
    err := toml.Unmarshal(b, &config)
    if err != nil {
        return nil, err
    }
    return &config, nil
}
