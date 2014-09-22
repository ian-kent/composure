package composure

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Composure map[string]*Composition

func NewComposure() *Composure {
	return &Composure{}
}

func (c *Composure) Get(name string) *Composition {
	cmp := (map[string]*Composition)(*c)[name]
	if cmp == nil {
		return nil
	}
	cmp.composure = c
	return cmp
}

func (c *Composure) NewComposition(name string) *Composition {
	cmp := &Composition{
		composure: c,
	}
	(map[string]*Composition)(*c)[name] = cmp
	return cmp
}

func ParseJSON(value []byte) (*Composure, error) {
	c := &Composure{}
	err := json.Unmarshal([]byte(value), c)
	if err != nil {
		return nil, err
	}
	for _, cm := range map[string]*Composition(*c) {
		cm.composure = c
		for _, cmp := range cm.Components {
			cmp.composition = cm
		}
		if cm.Template != nil {
			cm.Template.composition = cm
		}
		if cm.Preflight != nil {
			cm.Preflight.composition = cm
		}
		if cm.Postflight != nil {
			cm.Postflight.composition = cm
		}
	}
	return c, nil
}

func Load(file *os.File) (*Composure, error) {
	spec, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return ParseJSON(spec)
}
