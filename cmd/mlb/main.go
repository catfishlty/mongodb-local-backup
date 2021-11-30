package main

import (
    "github.com/alexflint/go-arg"
    log "github.com/sirupsen/logrus"
    "mongodb-local-backup/internal"
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
        log.Debugf("command: Start")
        if conf := getConfig(args.StartCmd.Config, args.StartCmd.Format, p, args.StartCmd.Daemon); conf != nil {
            if args.StartCmd.Daemon {
                runErr := internal.RunInDaemon(conf)
                if runErr != nil {
                    p.Fail(runErr.Error())
                    return
                }
            } else {
                log.Infof("mongodb-local-backup start")
                runErr := internal.Run(conf)
                if runErr != nil {
                    p.Fail(runErr.Error())
                    return
                }
                log.Infof("mongodb-local-backup stop")
            }
        }
    default:
        p.Fail("command not found")
    }
}

func getConfig(file, format string, p *arg.Parser, checkCron bool) *internal.Config {
    conf, err := internal.ReadConfig(file, format)
    if err != nil {
        p.Fail(err.Error())
        return nil
    }
    err = internal.CheckConfig(conf, checkCron)
    if err != nil {
        p.Fail(err.Error())
        return nil
    }
    return conf
}
