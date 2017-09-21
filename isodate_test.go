package golib

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestISODate(t *testing.T) {
	Convey("ISODate should work correctly", t, func() {

		Convey("Test now", func() {
			dt := Now()
			So(dt.Unix(), ShouldEqual, time.Now().Unix())
		})

		Convey("Test unix", func() {
			dt := Unix(1505907473, 0)
			So(dt.Unix(), ShouldEqual, 1505907473)
		})

		Convey("Test string", func() {
			dt := Unix(1505907473, 0)
			So(dt.String(), ShouldEqual, "2017-09-20T14:37:53.000+0300")
		})

		Convey("Test marshal json", func() {
			dt := Unix(1505907473, 0)
			b, _ := dt.MarshalJSON()
			So(string(b), ShouldEqual, "\"2017-09-20T14:37:53.000+0300\"")
		})

		Convey("Test marshal text", func() {
			dt := Unix(1505907473, 0)
			b, _ := dt.MarshalText()
			So(string(b), ShouldEqual, "2017-09-20T14:37:53.000+0300")
		})

		Convey("Test marshal bson", func() {
			dt := Unix(1505907473, 0)
			b, _ := dt.GetBSON()
			So(b.(time.Time), ShouldEqual, "2017-09-20 14:37:53 +0300 +03")
		})
	})

	Convey("ISODate  test JSON", t, func() {
		Convey("Test MarshalJson", func() {
			dt := Unix(1505907473, 0)
			b, _ := dt.MarshalJSON()
			So(string(b), ShouldEqual, "\"2017-09-20T14:37:53.000+0300\"")
		})

		Convey("Test UnmarshalJson", func() {
			dt := ISODate{}
			dt.UnmarshalJSON([]byte("2017-09-21T17:55:12.222Z"))
			So(dt, ShouldEqual, "\"2017-09-20T14:37:53.000+0300\"")
		})
	})

	Convey("ISODate ParseTimestamp", t, func() {
		Convey("Test ParseTimestamp", func() {
			b, _ := ParseTimestamp("2017-09-20T14:37:53.000+0300")
			So(b, ShouldEqual, "2017-09-20T14:37:53.000+0300")
		})
	})
}
