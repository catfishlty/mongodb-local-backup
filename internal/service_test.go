package internal

import (
	"fmt"
	. "github.com/agiledragon/gomonkey/v2"
	"github.com/agiledragon/gomonkey/v2/test/fake"
	"github.com/catfishlty/mongodb-local-backup/test"
	"github.com/commander-cli/cmd"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"io/fs"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestArgs_Version(t *testing.T) {
	Convey("TestArgs_Version", t, func() {
		arg := &Args{}
		So(arg.Version(), ShouldEqual, fmt.Sprintf("%s-%s-%s", Name, Version, runtime.GOOS))
	})
}

func TestGenTempFile(t *testing.T) {
	Convey("TestGetTempFile", t, func() {
		So(GenTempFile("json"), ShouldEqual, TempDir+BackupFileName+".json")
		So(GenTempFile("csv"), ShouldEqual, TempDir+BackupFileName+".csv")
	})
}

func TestGetTag(t *testing.T) {
	Convey("TestGetTag", t, func() {
		config := Config{
			Host: "127.0.0.1",
			Port: 80,
		}
		args := getArgs(&config)
		So(len(args), ShouldEqual, 2)
		So(args[0], ShouldEqual, "--host 127.0.0.1")
		So(args[1], ShouldEqual, "--port 80")
	})
}

func TestGetBackupFilename(t *testing.T) {
	prefix := "pre"
	prefixDefault := "mongodb-local-backup"
	db := "db"
	collection := "c"
	postfix := "post"
	Convey("TestGetBackupFilename", t, func() {
		Convey("test1", func() {
			s := GenBackupFilename(prefix, db, collection, postfix)
			So(strings.Index(s, prefix), ShouldEqual, 0)
			So(strings.Index(s, fmt.Sprintf("%s-%s-%s-", prefix, db, collection)), ShouldEqual, 0)
			So(strings.Index(s, postfix)+len(postfix), ShouldEqual, len(s))
		})
		Convey("test2", func() {
			s := GenBackupFilename("", db, collection, postfix)
			So(strings.Index(s, prefixDefault), ShouldEqual, 0)
			So(strings.Index(s, fmt.Sprintf("%s-%s-%s-", prefixDefault, db, collection)), ShouldEqual, 0)
			So(strings.Index(s, postfix)+len(postfix), ShouldEqual, len(s))
		})
	})
}

func getConfig() *Config {
	return &Config{
		Mongo: test.MongoexportPath,
		Host:  "127.0.0.1",
		Port:  27017,
		Target: []MongoTarget{
			{
				Db:         "test1",
				Collection: []string{"test11"},
			},
			{
				Db:         "test2",
				Collection: []string{"test21", "test22"},
			},
		},
		Type:   "json",
		Output: "mlb",
		Cron:   "*/1 * * * *",
		Prefix: "mlb",
	}
}

func TestRun(t *testing.T) {
	conf := getConfig()
	Convey("TestRun", t, func() {
		Convey("test1", func() {
			command := &cmd.Command{}
			patches := ApplyFunc(cmd.NewCommand, func(_ string, _ ...func(*cmd.Command)) *cmd.Command {
				return command
			})
			patches.ApplyMethodSeq(reflect.TypeOf(command), "Execute", []OutputCell{
				{Values: Params{nil}},
				{Values: Params{fake.ErrActual}},
				{Values: Params{nil}},
			})
			patches.ApplyMethodSeq(reflect.TypeOf(command), "ExitCode", []OutputCell{
				{Values: Params{0}},
				{Values: Params{0}},
				{Values: Params{2}},
			})
			patches.ApplyMethod(reflect.TypeOf(command), "Combined", func(*cmd.Command) string {
				return "test"
			})
			patches.ApplyFunc(ConvertByteToString, func([]byte, Charset) string {
				return "test"
			})
			patches.ApplyFuncSeq(FileTrans, []OutputCell{
				{Values: Params{nil}},
				{Values: Params{fake.ErrActual}},
			})
			defer patches.Reset()
			err := Run(conf)
			So(err, ShouldBeError)
		})
		Convey("test2", func() {
			command := &cmd.Command{}
			patches := ApplyFunc(cmd.NewCommand, func(_ string, _ ...func(*cmd.Command)) *cmd.Command {
				return command
			})
			patches.ApplyMethod(reflect.TypeOf(command), "Execute", func(*cmd.Command) error {
				return nil
			})
			patches.ApplyMethod(reflect.TypeOf(command), "ExitCode", func(*cmd.Command) int {
				return 0
			})
			patches.ApplyMethod(reflect.TypeOf(command), "Combined", func(*cmd.Command) string {
				return "test"
			})
			patches.ApplyFunc(ConvertByteToString, func([]byte, Charset) string {
				return "test"
			})
			patches.ApplyFunc(FileTrans, func(string, string) error { return nil })
			defer patches.Reset()
			err := Run(conf)
			So(err, ShouldBeNil)
		})
	})
}

func TestRunInDaemon(t *testing.T) {
	conf := getConfig()
	s := &gocron.Scheduler{}
	Convey("TestRunInDaemon", t, func() {
		Convey("test1", func() {
			patches :=
				ApplyFunc(os.OpenFile, func(string, int, fs.FileMode) (*os.File, error) {
					return nil, nil
				})
			patches.ApplyFunc(log.SetOutput, func(io.Writer) {})
			patches.ApplyMethod(reflect.TypeOf(s), "Do", func(*gocron.Scheduler, interface{}, ...interface{}) (*gocron.Job, error) {
				return nil, fake.ErrActual
			})
			defer patches.Reset()
			err := RunInDaemon(conf)
			So(err, ShouldEqual, fake.ErrActual)
		})
		Convey("test2", func() {
			confDup := getConfig()
			patches := ApplyFunc(os.OpenFile, func(string, int, fs.FileMode) (*os.File, error) {
				return nil, fake.ErrActual
			})
			patches.ApplyMethod(reflect.TypeOf(s), "Do", func(*gocron.Scheduler, interface{}, ...interface{}) (*gocron.Job, error) {
				return nil, fake.ErrActual
			})
			defer patches.Reset()
			err := RunInDaemon(confDup)
			So(err, ShouldEqual, fake.ErrActual)
		})
	})
}
