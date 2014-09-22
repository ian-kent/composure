package composure

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewComposure(t *testing.T) {
	Convey("NewComposure returns a new composure", t, func() {
		cmp := NewComposure()
		So(cmp, ShouldNotBeNil)
	})
}

func TestNewComposition(t *testing.T) {
	Convey("NewComposition should return a new composition", t, func() {
		cmp := NewComposure()

		So(map[string]*Composition(*cmp)["Foo"], ShouldBeNil)

		c := cmp.NewComposition("Foo")
		So(c, ShouldNotBeNil)
		So(map[string]*Composition(*cmp)["Foo"], ShouldNotBeNil)
	})
}

func TestGetComposition(t *testing.T) {
	Convey("Get should return a composition", t, func() {
		cmp := NewComposure()

		So(cmp.Get("Foo"), ShouldBeNil)

		c := cmp.NewComposition("Foo")
		So(cmp.Get("Foo"), ShouldNotBeNil)
		So(cmp.Get("Foo"), ShouldEqual, c)
	})
}

func TestParseJSON(t *testing.T) {
	Convey("ParseJSON should return an error for invalid JSON", t, func() {
		json := `Definitely invalid JSON`
		c, err := ParseJSON([]byte(json))

		So(c, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "invalid character 'D' looking for beginning of value")
	})
	Convey("ParseJSON should return a composure for valid JSON", t, func() {
		json := `{}`
		c, err := ParseJSON([]byte(json))

		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)
	})
	Convey("ParseJSON should create composition to composure reference", t, func() {
		json := `{"/":{}}`
		c, err := ParseJSON([]byte(json))

		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)
		So(c.Get("/").composure, ShouldEqual, c)
	})
}

func TestLoad(t *testing.T) {
	SkipConvey("Loading a non-existant file returns an error", t, func() {
		// TODO need a way to test file loading failures, maybe
		c, err := Load(&os.File{})

		So(c, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "EGGS")
	})
	Convey("Loading a file returns a composure", t, func() {
		f, err := os.Open("spec.json")
		if err == nil {
			c, err := Load(f)
			f.Close()

			So(c, ShouldNotBeNil)
			So(err, ShouldBeNil)
		} else {
			SkipSo(nil, ShouldEqual, nil)
		}
	})
}
