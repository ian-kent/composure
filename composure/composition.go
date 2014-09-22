package composure

import (
	"log"
	"net/http"
	"sync"
	"text/template"
)

type Composition struct {
	Methods    []string              `json:"methods"`
	Template   *Component            `json:"template"`
	Preflight  *Component            `json:"preflight"`
	Postflight *Component            `json:"postflight"`
	Components map[string]*Component `json:"components"`
	composure  *Composure
}

func (c *Composition) RenderFor(w http.ResponseWriter, req *http.Request) error {
	if c.Preflight != nil {
		// we don't care about the body content
		_, err := c.Preflight.Execute(w, req)
		if err != nil {
			return err
		}
	}

	tmpl, err := c.Template.Execute(w, req)
	if err != nil {
		return err
	}

	t, err := template.New("Text").Parse(string(tmpl))
	if err != nil {
		return err
	}

	args := make(map[string]string)
	if c.Components != nil {
		// FIXME scoping of wg seems wrong?
		var wg sync.WaitGroup
		for n, p := range c.Components {
			wg.Add(1)
			go func(n string, p *Component) {
				defer wg.Done()
				b, err := p.Execute(w, req)
				if err != nil {
					// FIXME better error handling?
					log.Fatal(err)
				}
				args[n] = string(b)
			}(n, p)
		}
		wg.Wait()
	}

	err = t.Execute(w, args)

	if err != nil {
		return err
	}

	if c.Postflight != nil {
		// we don't care about the body content
		_, err := c.Postflight.Execute(w, req)
		if err != nil {
			return err
		}
	}

	return nil
}
