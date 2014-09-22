package composure

import (
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetResource(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)

	c := NewComposure()
	cmp := c.NewComposition("Foo")
	cmp.Template = cmp.NewComponent("Text", []interface{}{}, []interface{}{})
	cmp.Template.Value = "Bar"

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Bar"))
	}))
	go http.ListenAndServe(":15677", nil)

	Convey("GetResource returns error for unknown type", t, func() {
		comp := cmp.NewComponent("Foo", []interface{}{}, []interface{}{})
		res, err := comp.GetResource(req)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "Unrecognised resource type: Foo")
		So(res, ShouldBeNil)
	})

	Convey("GetResource returns correct text resource", t, func() {
		comp := cmp.NewComponent("Text", []interface{}{}, []interface{}{})
		comp.Value = "Foo"
		res, err := comp.GetResource(req)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)

		bytes, err := ioutil.ReadAll(res.Reader)
		So(err, ShouldBeNil)
		So(string(bytes), ShouldEqual, "Foo")
	})

	Convey("GetResource returns correct URL resource", t, func() {
		comp := cmp.NewComponent("URL", []interface{}{"http://localhost:15677"}, []interface{}{})
		res, err := comp.GetResource(req)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)

		bytes, err := ioutil.ReadAll(res.Reader)
		So(err, ShouldBeNil)
		So(string(bytes), ShouldEqual, "Bar")
	})

	Convey("GetResource returns correct composition resource", t, func() {
		comp := cmp.NewComponent("Composition", []interface{}{}, []interface{}{})
		comp.Name = "Foo"
		res, err := comp.GetResource(req)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)

		bytes, err := ioutil.ReadAll(res.Reader)
		So(err, ShouldBeNil)
		So(string(bytes), ShouldEqual, "Bar")
	})

	Convey("getURLResource should return an error for missing params", t, func() {
		comp := cmp.NewComponent("URL", []interface{}{}, []interface{}{})
		res, err := comp.GetResource(req)
		So(res, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "Missing parameter for URL resource")

		comp = cmp.NewComponent("URL", []interface{}{map[string]interface{}{}}, []interface{}{})
		res, err = comp.GetResource(req)
		So(res, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "Missing URL parameter for complex URL resource")
	})

	Convey("getURLResource should return a result for complex URL type", t, func() {
		comp := cmp.NewComponent("URL", []interface{}{map[string]interface{}{"URL": "http://localhost:15677"}}, []interface{}{})
		res, err := comp.GetResource(req)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)

		bytes, err := ioutil.ReadAll(res.Reader)
		So(err, ShouldBeNil)
		So(string(bytes), ShouldEqual, "Bar")
	})

	Convey("getURLResource should return an error for an invalid URL", t, func() {
		comp := cmp.NewComponent("URL", []interface{}{"foobar://cant.exist"}, []interface{}{})
		res, err := comp.GetResource(req)
		So(res, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "Get foobar://cant.exist: unsupported protocol scheme \"foobar\"")
	})

	Convey("getCompositionResource should return an error for missing params", t, func() {
		comp := cmp.NewComponent("Composition", []interface{}{}, []interface{}{})
		res, err := comp.GetResource(req)
		So(res, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "Composition not found")
	})
}
