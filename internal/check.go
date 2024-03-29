package internal

import (
	"errors"
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/asaskevich/govalidator"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gorhill/cronexpr"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"os"
	"regexp"
	"strings"
)

// CheckArgsConfigFormat check config format is valid
func CheckArgsConfigFormat(p *arg.Parser, args Args) {
	if args.StartCmd != nil {
		if !slices.Contains([]string{"json", "yaml", "toml"}, args.StartCmd.Format) {
			errMsg := fmt.Sprintf("config file format '%s' is unsupported", args.StartCmd.Format)
			p.Fail(errMsg)
			log.Errorln(errMsg)
			os.Exit(1)
		}
	}
}

// CheckArgsLogLevel check log level is valid
func CheckArgsLogLevel(p *arg.Parser, args Args) {
	if args.StartCmd != nil {
		if !slices.Contains([]string{"panic", "fatal", "error", "warn", "warning", "info", "debug", "trace"}, args.StartCmd.LogLevel) {
			errMsg := fmt.Sprintf("log level '%s' is unsupported", args.StartCmd.LogLevel)
			p.Fail(errMsg)
			log.Errorln(errMsg)
			os.Exit(1)
		}
	}
}

// ValidFilePath custom validator to validate filepath for ozzo-validation
func ValidFilePath(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return errors.New("must be a valid file path")
	}
	ok, t := govalidator.IsFilePath(s)
	if ok {
		return nil
	}
	var osType string
	switch t {
	case govalidator.Win:
		osType = "windows"
	case govalidator.Unix:
		osType = "Unix"
	default:
		osType = ""
	}
	return errors.New("must be a valid " + osType + " file path")
}

// ValidCronExpression custom validator to validate cron expression for ozzo-validation
func ValidCronExpression(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return errors.New("must be a valid cron expression")
	}
	_, err := cronexpr.Parse(s)
	if err != nil {
		return errors.New("must be a valid cron expression: " + err.Error())
	}
	fields := strings.Split(s, " ")
	if len(fields) != 5 {
		return fmt.Errorf("fields must 5, not %d", len(fields))
	}
	return nil
}

// ValidDir custom validator to validate directory path for ozzo-validation
func ValidDir(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return errors.New("must be a valid dir path")
	}
	if !Exist(s) {
		err := Mkdir(s)
		if err != nil {
			return errors.New("must be a valid dir path: " + err.Error())
		}
	}
	if !IsDir(s) {
		return errors.New("'" + s + "' is not a valid dir path")
	}
	return nil
}

// ValidTarget custom validator to validate export db&collection in MongoDB for ozzo-validation
func ValidTarget(value interface{}) error {
	targets, ok := value.([]MongoTarget)
	if !ok {
		return errors.New("must be target list")
	}
	if len(targets) <= 0 {
		return errors.New("must greater than 0")
	}
	for i := 0; i < len(targets); i++ {
		target := targets[i]
		if target.Db == "" {
			return fmt.Errorf("db[%d] must not be empty", i)
		}
		if target.Collection == nil || len(target.Collection) <= 0 {
			return fmt.Errorf("db[%d] collection must not be nil and greater than 0", i)
		}
		for j := 0; j < len(target.Collection); j++ {
			if target.Collection[j] == "" {
				return fmt.Errorf("db[%d].collection[%d] must not be empty", i, j)
			}
		}
	}
	return nil
}

// CheckConfig call to validate all fields by ozzo-validation
func CheckConfig(conf *Config, cron bool) error {
	validList := make([]*validation.FieldRules, 0)
	validList = append(validList, validation.Field(&conf.Mongo, validation.Required, validation.By(ValidFilePath)))
	validList = append(validList, validation.Field(&conf.Host, validation.Required, is.Host))
	validList = append(validList, validation.Field(&conf.Port, validation.Required, validation.When(conf.Port > 0 && conf.Port < 65536)))
	validList = append(validList, validation.Field(&conf.Target, validation.Required, validation.By(ValidTarget)))
	validList = append(validList, validation.Field(&conf.Type, validation.Required, validation.In("json", "csv")))
	validList = append(validList, validation.Field(&conf.Output, validation.Required, validation.By(ValidFilePath), validation.By(ValidDir)))
	if cron {
		validList = append(validList, validation.Field(&conf.Cron, validation.Required, validation.By(ValidCronExpression)))
	}
	validList = append(validList, validation.Field(&conf.Prefix, validation.Match(regexp.MustCompile("[a-zA-Z0-9_-]+")).Error("the prefix must be character in 'a-zA-Z0-9_-'")))
	return validation.ValidateStruct(conf, validList...)
}
