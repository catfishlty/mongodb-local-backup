package main

import (
	"github.com/alexflint/go-arg"
	"github.com/catfishlty/mongodb-local-backup/internal"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

func main() {
	var args internal.Args
	p := arg.MustParse(&args)
	err := internal.CheckArgs(args)
	if err != nil {
		p.Fail(err.Error())
		return
	}

	switch {
	case args.StartCmd != nil:
		setLogLevel(p, args.StartCmd.LogLevel)
		log.Debugf("command: Start")
		conf := getEnvConfig(p, args.StartCmd.Daemon)
		if conf == nil {
			conf = getFileConfig(args.StartCmd.Config, args.StartCmd.Format, p, args.StartCmd.Daemon)
		}
		if conf == nil {
			p.Fail("missing config env variables or missing config file")
			return
		}
		if args.StartCmd.Daemon {
			exec(conf, p, internal.RunInDaemon)
		} else {
			exec(conf, p, internal.Run)
		}
	default:
		p.Fail("command not found")
	}
}

func setLogLevel(p *arg.Parser, levelString string) {
	level, err := log.ParseLevel(levelString)
	if err != nil {
		p.Fail("log level is not valid")
		os.Exit(1)
	}
	log.SetLevel(level)
}

func getEnvConfig(p *arg.Parser, daemon bool) *internal.Config {
	port, err := strconv.Atoi(os.Getenv(internal.MlbPort))
	if err != nil {
		p.Fail("port is not a number")
		os.Exit(1)
	}
	conf := &internal.Config{
		Mongo:    os.Getenv(internal.MlbMongoexport),
		Host:     os.Getenv(internal.MlbHost),
		Port:     port,
		Username: os.Getenv(internal.MlbUsername),
		Password: os.Getenv(internal.MlbPassword),
		Target:   getTargetConfig(os.Getenv(internal.MlbTarget)),
		Type:     os.Getenv(internal.MlbType),
		Output:   os.Getenv(internal.MlbOutput),
		Cron:     os.Getenv(internal.MlbCron),
		Prefix:   os.Getenv(internal.MlbPrefix),
		Log:      os.Getenv(internal.MlbLog),
	}
	err = internal.CheckConfig(conf, daemon)
	if err != nil {
		return nil
	}
	return conf
}

func getTargetConfig(configStr string) []internal.MongoTarget {
	target := make([]internal.MongoTarget, 0)
	dbColStrList := strings.Split(configStr, ";")
	for _, dbColStr := range dbColStrList {
		dbCol := strings.Split(dbColStr, "@")
		if len(dbCol) == 1 {
			target = append(target, internal.MongoTarget{
				Db: dbCol[0],
			})
		} else if len(dbCol) == 2 {
			target = append(target, internal.MongoTarget{
				Db:         dbCol[0],
				Collection: strings.Split(dbCol[1], ","),
			})
		}
	}
	return target
}

func getFileConfig(file, format string, p *arg.Parser, checkCron bool) *internal.Config {
	conf, err := internal.ReadConfig(file, format)
	if err != nil {
		p.Fail(err.Error())
		return nil
	}
	err = internal.CheckConfig(conf, checkCron)
	if err != nil {
		p.Fail(err.Error())

	}
	return conf
}

func exec(conf *internal.Config, p *arg.Parser, executor func(conf *internal.Config) error) {
	log.Infof("mongodb-local-backup start")
	runErr := executor(conf)
	if runErr != nil {
		p.Fail(runErr.Error())
		os.Exit(1)
	}
	log.Infof("mongodb-local-backup stop")
}
