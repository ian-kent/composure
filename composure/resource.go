package composure

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

type Resource struct {
	Reader   io.Reader
	Response *http.Response
}

func (c *Component) GetResource(req *http.Request) (*Resource, error) {
	switch c.Type {
	case "Text":
		return c.getTextResource(req, c.Value)
	case "HTTP":
		return c.getHTTPResource(req, c.Request)
	case "Composition":
		return c.getCompositionResource(req, c.Name)
	default:
		return nil, errors.New("Unrecognised resource type: " + c.Type)
	}
}

func (c *Component) getTextResource(req *http.Request, text string) (*Resource, error) {
	txt := strings.NewReader(text)

	return &Resource{txt, nil}, nil
}

func (c *Component) getHTTPResource(req *http.Request, params []interface{}) (*Resource, error) {
	if len(params) < 1 {
		return nil, errors.New("Missing parameter for HTTP resource")
	}

	var url string
	var args map[string]interface{}

	if s, ok := params[0].(string); ok {
		url = s
	} else {
		args = params[0].(map[string]interface{})
		if _, ok := args["URL"].(string); !ok {
			return nil, errors.New("Missing URL parameter for complex HTTP resource")
		}
		url = args["URL"].(string)
	}

	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if args != nil {
		if method, ok := args["Method"].(string); ok {
			r.Method = method
		}
		if hdrs, ok := args["Headers"].(map[string]interface{}); ok {
			if _, ok := hdrs["$INHERIT"]; req != nil && ok {
				for k, v := range req.Header {
					for _, v2 := range v {
						r.Header.Add(k, v2)
					}
				}
			}
			if remove, ok := hdrs["$REMOVE"].([]interface{}); ok {
				for _, k := range remove {
					r.Header.Del(k.(string))
				}
			}
			if copy, ok := hdrs["$COPY"].(map[string]interface{}); req != nil && ok {
				for k, v := range copy {
					for _, v2 := range req.Header[k] {
						r.Header.Add(v.(string), v2)
					}
				}
			}
			if set, ok := hdrs["$SET"].(map[string]interface{}); ok {
				for k, v := range set {
					r.Header.Set(k, v.(string))
				}
			}
			if add, ok := hdrs["$ADD"].(map[string]interface{}); ok {
				for k, v := range add {
					r.Header.Add(k, v.(string))
				}
			}
		}
	}

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	return &Resource{res.Body, res}, nil
}

func (c *Component) getCompositionResource(req *http.Request, name string) (*Resource, error) {
	r := c.composition.composure.Get(name)

	if r == nil {
		return nil, errors.New("Composition not found")
	}

	rr := httptest.NewRecorder()

	err := r.RenderFor(rr, req)
	if err != nil {
		return nil, err
	}

	// FIXME need to pass rr so we get headers etc
	return &Resource{rr.Body, nil}, nil
}
