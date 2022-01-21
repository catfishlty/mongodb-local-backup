package test

import (
    . "github.com/smartystreets/goconvey/convey"
    "io/fs"
    "testing"
    "time"
)

func TestMockFileInfo(t *testing.T) {
    f := MockFileInfo{}
    Convey("TestMockFileInfo", t, func() {
        Convey("name", func() {
            So(f.Name(), ShouldEqual, mockFileInfoName)
        })
        Convey("size", func() {
            So(f.Size(), ShouldEqual, mockFileInfoSize)
        })
        Convey("mode", func() {
            So(f.Mode(), ShouldEqual, fs.ModeDir)
        })
        Convey("mod time", func() {
            So(f.ModTime().UnixMilli(), ShouldBeLessThanOrEqualTo, time.Now().UnixMilli())
        })
        Convey("is dir", func() {
            So(f.IsDir(), ShouldBeTrue)
        })
        Convey("sys", func() {
            So(f.Sys(), ShouldEqual, mockFileInfoName)
        })
    })
}
