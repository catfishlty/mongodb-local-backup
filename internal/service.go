package internal

import (
	"fmt"
	"github.com/commander-cli/cmd"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"gopl.io/ch12/format"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"time"
)

// GenTempFile generate temp file path
func GenTempFile(format string) string {
	return fmt.Sprintf("%s%s.%s", TempDir, BackupFileName, format)
}

// Version gets the full version string
func (*Args) Version() string {
	return fmt.Sprintf("%s-%s-%s", Name, Version, runtime.GOOS)
}

// Run execute app with the config once
func Run(conf *Config) error {
	successList := make([]string, 0)
	failList := make([]string, 0)
	mongoSplitIdx := strings.LastIndex(conf.Mongo, string(os.PathSeparator))
	tempFile := GenTempFile(conf.Type)
	targets := conf.Target
	for _, target := range targets {
		for _, collection := range target.Collection {
			args := getArgs(conf)
			tempArgs := append(args, "--out "+tempFile)
			db := target.Db
			tempArgs = targetInject(tempArgs, db, collection)
			c := cmd.NewCommand(fmt.Sprintf("%s %s", conf.Mongo[mongoSplitIdx+1:], strings.Join(tempArgs, " ")), cmd.WithWorkingDir(conf.Mongo[:mongoSplitIdx]), cmd.WithEnvironmentVariables(cmd.EnvVars{"encoding": "utf-8"}))
			log.Debugf("mongo command: %s", c.Command)
			err := c.Execute()
			if err != nil || c.ExitCode() != 0 {
				errMsg := ""
				if err != nil {
					errMsg += err.Error()
				}
				if errMsg != "" {
					errMsg += " | "
				}
				errMsg += c.Combined()
				log.Errorf("mongo command execute error: %s", errMsg)
				failList = append(failList, fmt.Sprintf("%s:%s", db, collection))
				continue
			}
			log.Debugf("mongo response: %s", ConvertByteToString([]byte(c.Combined()), GB18030))
			err = FileTrans(tempFile, conf.Output+GenBackupFilename(conf.Prefix, db, collection, conf.Type))
			if err != nil {
				log.Errorf("file trans error: %v", err)
				failList = append(failList, fmt.Sprintf("%s:%s", db, collection))
				continue
			}
			successList = append(successList, fmt.Sprintf("%s:%s", db, collection))
		}
	}
	log.Infof("export success %d/%d", len(successList), len(successList)+len(failList))
	if len(successList) > 0 {
		log.Infof("success: [%s]", strings.Join(successList, ", "))
	}
	if len(failList) > 0 {
		log.Infof("fail: [%s]", strings.Join(failList, ", "))
	}
	if len(successList) != len(successList)+len(failList) {
		return fmt.Errorf("export err %d/%d", len(successList), len(successList)+len(failList))
	}
	return nil
}

// RunInDaemon execute app with the config by cron, and stop until shutdown the app
func RunInDaemon(conf *Config) error {
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
	s.StartAsync()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	s.Stop()
	return nil
}

func targetInject(args []string, db string, collection string) []string {
	args = append(args, fmt.Sprintf("--db %s", db))
	args = append(args, fmt.Sprintf("--collection %s", collection))
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
	return fmt.Sprintf("--%s %s", tag, val)
}

// GenBackupFilename generate backup file name
func GenBackupFilename(prefix, db, collection, postfix string) string {
	var p string
	if prefix != "" {
		p = prefix
	} else {
		p = "mongodb-local-backup"
	}
	return fmt.Sprintf("%s-%s-%s-%s.%s", p, db, collection, time.Now().Format("20060102150405-07MST"), postfix)
}
