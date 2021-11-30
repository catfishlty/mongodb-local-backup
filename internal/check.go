package internal

import (
    "errors"
    "fmt"
    "github.com/asaskevich/govalidator"
    validation "github.com/go-ozzo/ozzo-validation/v4"
    "github.com/go-ozzo/ozzo-validation/v4/is"
    "github.com/gorhill/cronexpr"
    "github.com/valyala/fasttemplate"
    "regexp"
    "strings"
)

func CheckArgs(args Args) error {
    m := map[string]bool{
        "json": true,
        "yaml": true,
        "toml": true,
    }
    switch {
    case args.StartCmd != nil:
        if !m[args.StartCmd.Format] {
            return errors.New("config file format '" + args.StartCmd.Format + "' is unsupported")
        }
    }
    return nil
}

func validFilePath(value interface{}) error {
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

func validCronExpression(value interface{}) error {
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
        return errors.New(fmt.Sprintf("fields must 5, not %d", len(fields)))
    }
    return nil
}

func validDir(value interface{}) error {
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

func validTarget(value interface{}) error {
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
            return errors.New(fasttemplate.ExecuteString(
                `db[{index}] must not be empty`, "{", "}",
                map[string]interface{}{"index": i},
            ))
        }
        if len(target.Collection) <= 0 {
            return errors.New(fasttemplate.ExecuteString(
                `db[{index}] collection must greater than 0`, "{", "}",
                map[string]interface{}{"index": i},
            ))
        }
        for j := 0; j < len(target.Collection); j++ {
            if target.Collection[j] == "" {
                return errors.New(fasttemplate.ExecuteString(
                    `db[{db}].collection[{col}] must not be empty`, "{", "}",
                    map[string]interface{}{"db": i, "col": j},
                ))
            }
        }
    }
    return nil
}

func CheckConfig(conf *Config, cron bool) error {
    validList := make([]*validation.FieldRules, 0)
    validList = append(validList, validation.Field(&conf.Mongo, validation.Required, validation.By(validFilePath)))
    validList = append(validList, validation.Field(&conf.Days, validation.When(conf.Days >= 0)))
    validList = append(validList, validation.Field(&conf.Host, validation.Required, is.Host))
    validList = append(validList, validation.Field(&conf.Port, validation.Required, validation.When(conf.Port > 0 && conf.Port < 65536)))
    validList = append(validList, validation.Field(&conf.Target, validation.Required, validation.By(validTarget)))
    validList = append(validList, validation.Field(&conf.Type, validation.Required, validation.In("json", "csv")))
    validList = append(validList, validation.Field(&conf.Output, validation.Required, validation.By(validFilePath), validation.By(validDir)))
    if cron {
        validList = append(validList, validation.Field(&conf.Cron, validation.Required, validation.By(validCronExpression)))
    }
    validList = append(validList, validation.Field(&conf.Prefix, validation.Match(regexp.MustCompile("[a-zA-Z0-9_-]+")).Error("the prefix must be characters")))
    return validation.ValidateStruct(conf, validList...)
}
