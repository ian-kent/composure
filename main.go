package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/pat"
	"github.com/ian-kent/composure/composure"
	"github.com/ian-kent/composure/context"
)

var comp *composure.Composure

type proxyHandler struct {
	Context *context.Context
}

func (ph proxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	cmp := comp.Get(ph.Context.Route)

	if cmp == nil {
		w.WriteHeader(404)
		return
	}

	cmp.RenderFor(w, req)
}

func listen() {
	gp := pat.New()

	for route, r := range *comp {
		if strings.HasPrefix(route, "/") {
			if len(r.Methods) > 0 {
				for _, m := range r.Methods {
					gp.Add(m, route, proxyHandler{&context.Context{Method: m, Route: route}})
				}
			} else {
				gp.Add("GET", route, proxyHandler{&context.Context{Method: "GET", Route: route}})
			}
		}
	}

	err := http.ListenAndServe(":9000", gp)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func main() {
	flag.Parse()

	filename := "spec.json"
	if len(flag.Args()) > 0 {
		filename = flag.Args()[0]
	}

	var spec *composure.Composure
	var err error

	if strings.HasPrefix(filename, "http://") ||
		strings.HasPrefix(filename, "https://") {
		res, err := http.Get(filename)
		if err != nil {
			log.Fatal(err)
		}

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		res.Body.Close()

		spec, err = composure.ParseJSON(b)
	} else {
		f, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}

		spec, err = composure.Load(f)
		f.Close()
	}

	if err != nil {
		log.Fatal(err)
	}

	if spec == nil {
		log.Fatal(errors.New("Error loading spec file"))
	}

	comp = spec

	log.Println("Loaded spec: " + filename)

	listen()
}
