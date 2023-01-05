package internal

import (
	"errors"
	. "github.com/agiledragon/gomonkey/v2"
	"github.com/agiledragon/gomonkey/v2/test/fake"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCheckArgs(t *testing.T) {
	Convey("TestCheckArgs", t, func() {
		Convey("test1", func() {
			args := Args{
				StartCmd: nil,
			}
			So(CheckArgs(nil, args), ShouldBeNil)
		})
		Convey("test2", func() {
			args := Args{
				StartCmd: &BaseCmd{
					Format: "json",
				},
			}
			So(CheckArgs(nil, args), ShouldBeNil)
		})
		Convey("test3", func() {
			args := Args{
				StartCmd: &BaseCmd{
					Format: "xxx",
				},
			}
			So(CheckArgs(nil, args), ShouldBeError)
		})
	})
}

func TestValidFilePath(t *testing.T) {
	Convey("TestValidFilePath", t, func() {
		Convey("test1", func() {
			So(ValidFilePath(33), ShouldNotEqual, nil)
		})
		Convey("test2", func() {
			So(ValidFilePath("C:\\Program Files\\"), ShouldEqual, nil)
		})
		Convey("test3", func() {
			So(ValidFilePath("/bin/bash"), ShouldEqual, nil)
		})
		Convey("test4", func() {
			So(ValidFilePath("C:\\:\\"), ShouldNotEqual, nil)
		})
	})
}

func TestValidCronExpression(t *testing.T) {
	Convey("TestValidCronExpression", t, func() {
		Convey("test1", func() {
			So(ValidCronExpression(33), ShouldBeError)
		})
		Convey("test2", func() {
			So(ValidCronExpression("test"), ShouldBeError)
		})
		Convey("test3", func() {
			So(ValidCronExpression("0 * * * * *"), ShouldBeError)
		})
		Convey("test4", func() {
			So(ValidCronExpression("*/1 * * * *"), ShouldBeNil)
		})
	})
}

func TestValidDir(t *testing.T) {
	Convey("TestValidDir", t, func() {
		Convey("test1", func() {
			So(ValidDir(33), ShouldBeError)
		})
		Convey("test2", func() {
			patches := ApplyFunc(Exist, func(_ string) bool {
				return false
			})
			patches.ApplyFunc(Mkdir, func(_ string) error {
				return errors.New("test")
			})
			defer patches.Reset()
			So(ValidDir("test"), ShouldBeError)
		})
		Convey("test3", func() {
			existPatch := ApplyFunc(Exist, func(_ string) bool {
				return false
			})
			defer existPatch.Reset()
			mkdirPatch := ApplyFunc(Mkdir, func(_ string) error {
				return nil
			})
			defer mkdirPatch.Reset()
			isDirPatch := ApplyFunc(IsDir, func(_ string) bool {
				return false
			})
			defer isDirPatch.Reset()
			So(ValidDir("test"), ShouldBeError)
		})
		Convey("test4", func() {
			existPatch := ApplyFunc(Exist, func(_ string) bool {
				return true
			})
			defer existPatch.Reset()
			isDirPatch := ApplyFunc(IsDir, func(_ string) bool {
				return false
			})
			defer isDirPatch.Reset()
			So(ValidDir("test"), ShouldBeError)
		})
		Convey("test5", func() {
			existPatch := ApplyFunc(Exist, func(_ string) bool {
				return true
			})
			defer existPatch.Reset()
			isDirPatch := ApplyFunc(IsDir, func(_ string) bool {
				return true
			})
			defer isDirPatch.Reset()
			So(ValidDir("test"), ShouldBeNil)
		})
	})
}

func TestValidTarget(t *testing.T) {
	Convey("TestValidTarget", t, func() {
		Convey("test1", func() {
			So(ValidTarget(33), ShouldBeError)
		})
		Convey("test2", func() {
			data := make([]MongoTarget, 0)
			So(ValidTarget(data), ShouldBeError)
		})
		Convey("test3", func() {
			data := make([]MongoTarget, 0)
			data = append(data, MongoTarget{
				Db:         "",
				Collection: nil,
			})
			So(ValidTarget(data), ShouldBeError)
		})
		Convey("test4", func() {
			data := make([]MongoTarget, 0)
			data = append(data, MongoTarget{
				Db:         "test",
				Collection: nil,
			})
			So(ValidTarget(data), ShouldBeError)
		})
		Convey("test5", func() {
			data := make([]MongoTarget, 0)
			data = append(data, MongoTarget{
				Db:         "test",
				Collection: []string{},
			})
			So(ValidTarget(data), ShouldBeError)
		})
		Convey("test6", func() {
			data := make([]MongoTarget, 0)
			data = append(data, MongoTarget{
				Db:         "test",
				Collection: []string{""},
			})
			So(ValidTarget(data), ShouldBeError)
		})
		Convey("test7", func() {
			data := make([]MongoTarget, 0)
			data = append(data, MongoTarget{
				Db:         "test1",
				Collection: []string{"test11"},
			})
			data = append(data, MongoTarget{
				Db:         "test2",
				Collection: []string{"test21", "test22"},
			})
			So(ValidTarget(data), ShouldBeNil)
		})
	})
}

func TestCheckConfig(t *testing.T) {
	conf := &Config{
		Mongo: "/usr/sbin/mongoexport",
		Host:  "127.0.0.1",
		Port:  27017,
		Target: []MongoTarget{
			{Db: "test", Collection: []string{"test1", "test2"}},
		},
		Type:   "json",
		Output: "/root/output",
		Cron:   "*/1 * * * *",
		Prefix: "mongo-local-backup",
		Log:    "/root/logs",
	}
	Convey("TestCheckConfig", t, func() {
		Convey("test1", func() {
			patches := ApplyFunc(Mkdir, func(string) error {
				return nil
			})
			patches.ApplyFunc(validation.ValidateStruct, func(interface{}, ...*validation.FieldRules) error {
				return nil
			})
			defer patches.Reset()
			So(CheckConfig(conf, false), ShouldBeNil)
			So(CheckConfig(conf, true), ShouldBeNil)
		})
		Convey("test2", func() {
			patches := ApplyFunc(Mkdir, func(string) error {
				return nil
			})
			patches.ApplyFunc(validation.ValidateStruct, func(interface{}, ...*validation.FieldRules) error {
				return fake.ErrActual
			})
			defer patches.Reset()
			So(CheckConfig(conf, false), ShouldBeError)
			So(CheckConfig(conf, false), ShouldEqual, fake.ErrActual)
			So(CheckConfig(conf, true), ShouldBeError)
			So(CheckConfig(conf, true), ShouldEqual, fake.ErrActual)
		})
	})
}
