package internal

import (
    "errors"
    "fmt"
    "github.com/commander-cli/cmd"
    "github.com/go-co-op/gocron"
    log "github.com/sirupsen/logrus"
    "github.com/valyala/fasttemplate"
    "gopl.io/ch12/format"
    "os"
    "os/signal"
    "reflect"
    "runtime"
    "strings"
    "syscall"
    "time"
)

func GetTempFile(format string) string {
    return fmt.Sprintf("%s%s.%s", TempDir, BackupFileName, format)
}

func (*Args) Version() string {
    return fmt.Sprintf("%s-%s-%s", Name, Version, runtime.GOOS)
}

func Run(conf *Config) error {
    successList := make([]string, 0)
    failList := make([]string, 0)
    mongoSplitIdx := strings.LastIndex(conf.Mongo, "\\")
    tempFile := GetTempFile(conf.Type)
    targets := conf.Target
    args := getArgs(conf)
    for i := 0; i < len(targets); i++ {
        tempArgs := append(args, "--out "+tempFile)
        for j := 0; j < len(targets[i].Collection); j++ {
            db := targets[i].Db
            collection := targets[i].Collection[j]
            tempArgs = targetInject(tempArgs, db, collection)
            c := cmd.NewCommand(fmt.Sprintf("%s %s", conf.Mongo[mongoSplitIdx+1:], strings.Join(tempArgs, " ")), cmd.WithWorkingDir(conf.Mongo[:mongoSplitIdx]), cmd.WithEnvironmentVariables(cmd.EnvVars{"encoding": "utf-8"}))
            log.Debugf("mongo command: %s", c.Command)
            _ = c.Execute()
            log.Debugf("mongo response: %s", ConvertByteToString([]byte(c.Combined()), GB18030))
            err := FileTrans(tempFile, conf.Output+GetBackupFilename(conf.Prefix, db, collection, conf.Type))
            if err != nil {
                log.Errorf("file trans error: %v", err)
                failList = append(failList, fmt.Sprintf("%s:%s", db, collection))
                continue
            }
            successList = append(successList, fmt.Sprintf("%s:%s", db, collection))
        }
    }
    if len(successList) != len(successList)+len(failList) {
        return errors.New(fasttemplate.ExecuteString(`export err {current}/{total}`, "{", "}", map[string]interface{}{
            "current": len(successList),
            "total":   len(successList) + len(failList),
        }))
    }
    log.Infof("export success %d/%d", len(successList), len(successList)+len(failList))
    if len(successList) > 0 {
        log.Infof("success: [%s]", strings.Join(successList, ", "))
    }
    if len(failList) > 0 {
        log.Infof("fail: [%s]", strings.Join(failList, ", "))
    }
    return nil
}

func RunInDaemon(conf *Config) error {
    initLogger(getLogPath(conf.Log))
    s := gocron.NewScheduler(time.Now().Location())
    _, err := s.Cron(conf.Cron).Do(func() {
        runErr := Run(conf)
        if runErr != nil {
            log.Errorf("run cron job fail: %s", runErr.Error())
            return
        }
    })
    if err != nil {
        log.Errorf("add cron job fail: %s", err.Error())
        return err
    }
    log.Infof("mongodb-local-backup[daemon] start")
    s.StartAsync()
    quit := make(chan os.Signal)
    signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
    <-quit
    s.Stop()
    log.Infof("mongodb-local-backup[daemon] stop")
    return nil
}

func targetInject(args []string, db string, collection string) []string {
    args = append(args, fasttemplate.ExecuteString(`--db {val}`, "{", "}", map[string]interface{}{"val": db}))
    args = append(args, fasttemplate.ExecuteString(`--collection {val}`, "{", "}", map[string]interface{}{"val": collection}))
    return args
}

func getArgs(conf *Config) []string {
    args := make([]string, 0)
    getAnnotation(*conf, func(tag, val string) {
        args = append(args, formatCmd(tag, val))
    })
    return args
}

func getAnnotation(e interface{}, exec func(tag, val string)) {
    t := reflect.TypeOf(e)
    v := reflect.ValueOf(e)
    for i := 0; i < t.NumField(); i++ {
        ft := t.Field(i)
        fv := v.Field(i)
        tagCmd := ft.Tag.Get(TagCmd)
        if tagCmd != "" {
            val := strings.Trim(format.Any(fv.Interface()), "\"")
            if val != "" {
                exec(tagCmd, val)
            }
        }
    }
}

func formatCmd(tag, val string) string {
    return fasttemplate.ExecuteString(`--{tag} {val}`, "{", "}",
        map[string]interface{}{
            "tag": tag,
            "val": val,
        })
}

func GetBackupFilename(prefix, db, collection, postfix string) string {
    var p string
    if prefix != "" {
        p = prefix
    } else {
        p = "mongodb-local-backup"
    }
    t := fasttemplate.New(`{prefix}-{db}-{collection}-{time}.{postfix}`, "{", "}")
    return t.ExecuteString(map[string]interface{}{
        "prefix":     p,
        "db":         db,
        "collection": collection,
        "time":       time.Now().Format("20060102150405999-07MST"),
        "postfix":    postfix,
    })
}

func getLogPath(logPath string) string {
    if logPath == "" {
        return "mongodb-local-backup.log"
    }
    return logPath
}

func initLogger(logFile string) {
    f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND, os.ModeAppend|os.ModePerm)
    if err != nil {
        log.Errorf("log file path error[%s]: %s", logFile, err.Error())
        return
    }
    log.SetOutput(f)
}
