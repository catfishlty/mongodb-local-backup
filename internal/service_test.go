package internal

import (
	"testing"
)

func TestGetTempFile(t *testing.T) {
    Convey("TestGetTempFile", t, func() {
        So(GetTempFile("json"), ShouldEqual, TempDir+BackupFileName+".json")
        So(GetTempFile("csv"), ShouldEqual, TempDir+BackupFileName+".csv")
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
