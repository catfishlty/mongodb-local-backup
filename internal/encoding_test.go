package internal

import (
	. "github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/text/encoding"
	"reflect"
	"testing"
)

func TestConvertByteToString(t *testing.T) {
	d := &encoding.Decoder{}
	s := "test"
	Convey("TestConvertByteToString", t, func() {
		Convey("gb18030 test", func() {
			patches := ApplyMethod(reflect.TypeOf(d), "Bytes", func(_ *encoding.Decoder, _ []byte) ([]byte, error) {
				return []byte(s), nil
			})
			defer patches.Reset()
			So(ConvertByteToString([]byte{0}, GB18030), ShouldEqual, s)
		})
		Convey("utf8 test", func() {
			So(ConvertByteToString([]byte(s), UFT8), ShouldEqual, s)
		})
		Convey("others test", func() {
			So(ConvertByteToString([]byte(s), "GBK"), ShouldEqual, s)
		})
	})
}
