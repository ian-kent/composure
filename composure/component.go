package composure

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Component struct {
	Type        string        `json:"type"`
	Request     []interface{} `json:"request"`
	Response    []interface{} `json:"response"`
	Value       string        `json:"value"`
	Name        string        `json:"name"`
	Method      string        `json:"method"`
	composition *Composition
}

func (c *Composition) NewComponent(name string, request []interface{}, response []interface{}) *Component {
	// FIXME feels like this should add it to c.Components, but that
	// breaks for Template and Pre-/Post-flight components
	cmp := &Component{
		Type:        name,
		Request:     request,
		Response:    response,
		composition: c,
	}
	return cmp
}

func (c *Component) Execute(w http.ResponseWriter, req *http.Request) ([]byte, error) {
	var b []byte

	// FIXME doesn't work for Response only URL types, ugh
	r, err := c.GetResource(req)
	if err != nil {
		log.Println(err.Error())
		//return b, err
	}
	c.Transform(w, req, c.Response, r)
	if err != nil {
		return b, err
	}
	return ioutil.ReadAll(io.Reader(r.Reader))
}

func (c *Component) Transform(w http.ResponseWriter, req *http.Request, response []interface{}, r *Resource) {
	var args map[string]interface{}

	if response != nil && len(response) > 0 {
		if a, ok := response[0].(map[string]interface{}); ok {
			args = a
		}
	}

	if args != nil {
		if hdrs, ok := args["Headers"].(map[string]interface{}); ok {
			if _, ok := hdrs["$INHERIT"]; req != nil && ok {
				for k, v := range r.Response.Header {
					for _, v2 := range v {
						w.Header().Add(k, v2)
					}
				}
			}
			if remove, ok := hdrs["$REMOVE"].([]interface{}); ok {
				for _, k := range remove {
					w.Header().Del(k.(string))
				}
			}
			if copy, ok := hdrs["$COPY"].(map[string]interface{}); req != nil && ok {
				for k, v := range copy {
					for _, v2 := range r.Response.Header[k] {
						w.Header().Add(v.(string), v2)
					}
				}
			}
			if set, ok := hdrs["$SET"].(map[string]interface{}); ok {
				for k, v := range set {
					w.Header().Set(k, v.(string))
				}
			}
			if add, ok := hdrs["$ADD"].(map[string]interface{}); ok {
				for k, v := range add {
					w.Header().Add(k, v.(string))
				}
			}
		}
	}
}
