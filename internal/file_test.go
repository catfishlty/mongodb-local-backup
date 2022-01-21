package internal

import (
    "encoding/json"
    "github.com/BurntSushi/toml"
    . "github.com/agiledragon/gomonkey/v2"
    "github.com/agiledragon/gomonkey/v2/test/fake"
    "github.com/catfishlty/mongodb-local-backup/test"
    . "github.com/smartystreets/goconvey/convey"
    "gopkg.in/yaml.v3"
    "io/fs"
    "io/ioutil"
    "os"
    "reflect"
    "testing"
)

func TestExist(t *testing.T) {
    Convey("TestExist", t, func() {
        Convey("test1", func() {
            patches := ApplyFunc(os.Stat, func(string) (os.FileInfo, error) {
                return nil, fake.ErrActual
            })
            patches.ApplyFunc(os.IsExist, func(_ error) bool {
                return true
            })
            defer patches.Reset()
            So(Exist("test"), ShouldBeTrue)
        })
        Convey("test2", func() {
            patches := ApplyFunc(os.Stat, func(string) (os.FileInfo, error) {
                return nil, fake.ErrActual
            })
            patches.ApplyFunc(os.IsExist, func(_ error) bool {
                return false
            })
            defer patches.Reset()
            So(Exist("test"), ShouldBeFalse)
        })
        Convey("test3", func() {
            patches := ApplyFunc(os.Stat, func(string) (os.FileInfo, error) {
                return nil, nil
            })
            defer patches.Reset()
            So(Exist("test"), ShouldBeTrue)
        })
    })
}

func TestIsDir(t *testing.T) {
    f := &test.MockFileInfo{}
    Convey("TestIsDir", t, func() {
        Convey("test1", func() {
            patches := ApplyFunc(os.Stat, func(string) (os.FileInfo, error) {
                return f, fake.ErrActual
            })
            defer patches.Reset()
            So(IsDir("test"), ShouldBeFalse)
        })
        Convey("test2", func() {
            patches := ApplyFunc(os.Stat, func(string) (os.FileInfo, error) {
                return f, nil
            })
            patches.ApplyMethod(reflect.TypeOf(f), "IsDir", func(_ *test.MockFileInfo) bool {
                return false
            })
            defer patches.Reset()
            So(IsDir("test"), ShouldBeFalse)
        })
        Convey("test3", func() {
            patches := ApplyFunc(os.Stat, func(string) (os.FileInfo, error) {
                return f, nil
            })
            patches.ApplyMethod(reflect.TypeOf(f), "IsDir", func(_ *test.MockFileInfo) bool {
                return true
            })
            defer patches.Reset()
            So(IsDir("test"), ShouldBeTrue)
        })
    })
}

func TestMkdir(t *testing.T) {
    Convey("TestMkdir", t, func() {
        Convey("test1", func() {
            patches := ApplyFunc(os.MkdirAll, func(_ string, _ fs.FileMode) error {
                return nil
            })
            defer patches.Reset()
            So(Mkdir("test"), ShouldBeNil)
        })
        Convey("test2", func() {
            patches := ApplyFunc(os.MkdirAll, func(_ string, _ fs.FileMode) error {
                return fake.ErrActual
            })
            defer patches.Reset()
            So(Mkdir("test"), ShouldBeError)
            So(Mkdir("test"), ShouldEqual, fake.ErrActual)
        })
    })
}

func TestIsFile(t *testing.T) {
    Convey("TestIsFile", t, func() {
        Convey("test1", func() {
            patches := ApplyFunc(IsDir, func(string) bool {
                return true
            })
            defer patches.Reset()
            So(IsFile("test"), ShouldBeFalse)
        })
        Convey("test2", func() {
            patches := ApplyFunc(IsDir, func(string) bool {
                return false
            })
            defer patches.Reset()
            So(IsFile("test"), ShouldBeTrue)
        })
    })
}

func TestFileTrans(t *testing.T) {
    Convey("TestFileTrans", t, func() {
        Convey("test1", func() {
            patches := ApplyFunc(ioutil.ReadFile, func(string) ([]byte, error) {
                return nil, nil
            })
            patches.ApplyFunc(ioutil.WriteFile, func(string, []byte, fs.FileMode) error {
                return nil
            })
            patches.ApplyFunc(os.Remove, func(string) error {
                return nil
            })
            defer patches.Reset()
            So(FileTrans("test", "test"), ShouldBeNil)
        })
        Convey("test2", func() {
            patches := ApplyFunc(ioutil.ReadFile, func(string) ([]byte, error) {
                return nil, nil
            })
            patches.ApplyFunc(ioutil.WriteFile, func(string, []byte, fs.FileMode) error {
                return nil
            })
            patches.ApplyFunc(os.Remove, func(string) error {
                return fake.ErrActual
            })
            defer patches.Reset()
            So(FileTrans("test", "test"), ShouldBeError)
            So(FileTrans("test", "test"), ShouldEqual, fake.ErrActual)
        })
        Convey("test3", func() {
            patches := ApplyFunc(ioutil.ReadFile, func(string) ([]byte, error) {
                return nil, nil
            })
            patches.ApplyFunc(ioutil.WriteFile, func(string, []byte, fs.FileMode) error {
                return fake.ErrActual
            })
            defer patches.Reset()
            So(FileTrans("test", "test"), ShouldBeError)
            So(FileTrans("test", "test"), ShouldEqual, fake.ErrActual)
        })
        Convey("test4", func() {
            patches := ApplyFunc(ioutil.ReadFile, func(string) ([]byte, error) {
                return nil, fake.ErrActual
            })
            defer patches.Reset()
            So(FileTrans("test", "test"), ShouldBeError)
            So(FileTrans("test", "test"), ShouldEqual, fake.ErrActual)
        })
    })
}

func TestReadConfig(t *testing.T) {
    Convey("TestReadConfig", t, func() {
        Convey("test1", func() {
            patches := ApplyFunc(ioutil.ReadFile, func(string) ([]byte, error) {
                return nil, fake.ErrActual
            })
            defer patches.Reset()
            c, err := ReadConfig("test", "test")
            So(c, ShouldBeNil)
            So(err, ShouldBeError)
            So(err, ShouldEqual, fake.ErrActual)
        })
        Convey("test2", func() {
            patches := ApplyFunc(ioutil.ReadFile, func(string) ([]byte, error) {
                return nil, nil
            })
            patches.ApplyFunc(ReadJson, func([]byte) (*Config, error) {
                return &Config{Type: "json"}, nil
            })
            patches.ApplyFunc(ReadYaml, func([]byte) (*Config, error) {
                return &Config{Type: "yaml"}, nil
            })
            patches.ApplyFunc(ReadToml, func([]byte) (*Config, error) {
                return &Config{Type: "toml"}, nil
            })
            defer patches.Reset()
            c, err := ReadConfig("test", "json")
            So(c, ShouldNotBeNil)
            So(c.Type, ShouldEqual, "json")
            So(err, ShouldBeNil)
            c, err = ReadConfig("test", "yaml")
            So(c, ShouldNotBeNil)
            So(c.Type, ShouldEqual, "yaml")
            So(err, ShouldBeNil)
            c, err = ReadConfig("test", "toml")
            So(c, ShouldNotBeNil)
            So(c.Type, ShouldEqual, "toml")
            So(err, ShouldBeNil)
            c, err = ReadConfig("test", "test")
            So(c, ShouldBeNil)
            So(err, ShouldBeError)
        })
    })
}

func TestReadJson(t *testing.T) {
    Convey("TestReadJson", t, func() {
        Convey("test1", func() {
            patches := ApplyFunc(json.Unmarshal, func([]byte, interface{}) error {
                return fake.ErrActual
            })
            defer patches.Reset()
            c, err := ReadJson([]byte{})
            So(err, ShouldBeError)
            So(err, ShouldEqual, fake.ErrActual)
            So(c, ShouldBeNil)
        })
        Convey("test2", func() {
            patches := ApplyFunc(json.Unmarshal, func([]byte, interface{}) error {
                return nil
            })
            defer patches.Reset()
            c, err := ReadJson([]byte{})
            So(err, ShouldBeNil)
            So(c, ShouldBeNil)
        })
    })
}
func TestReadYaml(t *testing.T) {
    Convey("TestReadYaml", t, func() {
        Convey("test1", func() {
            patches := ApplyFunc(yaml.Unmarshal, func([]byte, interface{}) error {
                return fake.ErrActual
            })
            defer patches.Reset()
            c, err := ReadYaml([]byte{})
            So(err, ShouldBeError)
            So(err, ShouldEqual, fake.ErrActual)
            So(c, ShouldBeNil)
        })
        Convey("test2", func() {
            patches := ApplyFunc(yaml.Unmarshal, func([]byte, interface{}) error {
                return nil
            })
            defer patches.Reset()
            _, err := ReadYaml([]byte{})
            So(err, ShouldBeNil)
        })
    })
}
func TestReadToml(t *testing.T) {
    Convey("TestReadToml", t, func() {
        Convey("test1", func() {
            patches := ApplyFunc(toml.Unmarshal, func([]byte, interface{}) error {
                return fake.ErrActual
            })
            defer patches.Reset()
            c, err := ReadToml([]byte{})
            So(err, ShouldBeError)
            So(err, ShouldEqual, fake.ErrActual)
            So(c, ShouldBeNil)
        })
        Convey("test2", func() {
            patches := ApplyFunc(toml.Unmarshal, func([]byte, interface{}) error {
                return nil
            })
            defer patches.Reset()
            _, err := ReadToml([]byte{})
            So(err, ShouldBeNil)
        })
    })
}
