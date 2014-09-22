Composure
=========

[![Build Status](https://travis-ci.org/ian-kent/composure.svg?branch=master)](https://travis-ci.org/ian-kent/composure)

Composure is a lightweight server-side UI composition framework for building
web applications with microservices.

* Compose distinct UIs and APIs into one service
* Transform HTTP requests and responses
* Define compositions in JSON
* Load templates and components from remote services
* Asynchronous rendering - as fast as your slowest component
* Integrate legacy services to avoid component rewrites

### Quick start

* Run `composure examples/demo/spec.json` to load the [Composure demo](examples/demo)
* Open the Composure in your browser: http://localhost:9000
* Run it with [Docker](https://www.docker.com/) - see [the example Dockerfile](examples/demo/Dockerfile)

:bulb: You can specify a URL to load, e.g. `composure http://your/spec.json`

:warning: Composure is an experimental product, not suitable for production environments.

### Composure web application example

This is an example of using Composure to combine web components written in multiple
languages into one user interface.

<img src="https://docs.google.com/drawings/d/1tppdtr7noODf_1qJYOORdFmnmSaRh4mwZZxMiTkDBqQ/pub?w=929&amp;h=682">

#### Example JSON specification

This example specification implements the diagram above.

It uses the following configuration:

* **Preflight**

  The `Authorization` header is mapped to `X-Composure-Authorization`.

  The request is synchronously sent to an authentication server, which validates
  `X-Composure-Authorization` and responds with `X-Composure-Identity`.

  `X-Composure-Identity` is mapped to `X-Identity` for downstream requests.

* **Components**

  Requests to render each component are sent asynchronously.

  Rendering services use `X-Identity` to identify the user.

* **Template**

  Renders the template synchronously using the component output.

* **Postflight**

  Transforms the final response synchronously, to remove the `X-Identity`
  header and set the `Server` header to `Composure`.

```json
{
  "/": {
    "Template": {
      "Type": "URL",
      "Request": [ "http://templates.foobar/index.tmpl" ]
    },
    "Components": {
      "Header": {
        "Type": "URL",
        "Request": [ "http://header.foobar" ]
      },
      "Navbar": {
        "Type": "URL",
        "Request": [ "http://navbar.foobar" ]
      },
      "Content": {
        "Type": "URL",
        "Request": [ "http://newsfeed.foobar" ]
      },
      "Footer": {
        "Type": "URL",
        "Request": [ "http://footer.foobar" ]
      },
    },
    "Preflight": {
      "Type": "URL",
      "Request": [{
        "URL": "http://authentication.foobar/preflight",
        "Method": "GET",
        "Headers": {
          "$COPY": {
            "Authorization": "X-Composure-Authorization"
          }
        }
      }],
      "Response": [{
        "Headers": {
          "$REMOVE": [ "X-Composure-Authorization" ],
          "$COPY": {
            "X-Composure-Identity": "X-Identity"
          }
        }
      }]
    },
    "Postflight": {
      "Type": "URL",
      "Response": [{
        "Headers": {
          "$SET": {
            "Server": "Composure"
          },
          "$REMOVE": [ "X-Identity" ]
        }
      }]
    }
  }
}
```

### JSON specification

Composition names beginning with `/` are routes served by `composure`.

Routes are defined using [Gorilla Pat](https://github.com/gorilla/pat).

Other compositions are reserved for internal use, e.g. to build reusable components.

* Preflight is executed synchronously before rendering takes place
* Components are rendered asynchronously after Preflight
* Template is rendered synchronously after all Components are rendered
* Postflight is executed synchronously after all rendering is complete

On a route composition, you can optionally specify HTTP methods:

```json
"/": {
  "Methods": [ "GET" ]
}
```

If no methods are specified, GET is registered by default.

#### Composition

Each child object is a "composition":

```json
"Navbar": {
  "Template": {
    "Type": "Text",
    "Value": "<p>Nav bar!</p>"
  }
}
```

A composition consists of a template and, optionally, some components.

Templates are essentially Go `html/template` templates, and components
are injected into the template using `html/template` syntax, e.g.

```
{{ .Navbar }}
{{ .Content }}
```

This means you can use logic in your templates to control rendering of
remote components as if they were local variables, e.g.

```
{{ if .Admin }}
  <div class="admin_nav">{{ .Admin }}</div>
{{ end }}
```

The Template, Preflight and Postflight parameters are all Components,
and can use any of the parameters defined below.

#### Component

A component is a JSON object specifying a `Type`:

```json
"Navbar": {
  "Type": "Composition",
  "Name": "Navbar"
}
```

Additional items are passed to the component:

* URL - loads the specified URL (see below)
* Text - uses the raw text in `Value`
* Composition - renders composition named in `Name` from this composure

#### HTTP requests and responses

The URL component type supports simple requests and complex requests:

The simple form takes a single string argument, generating a GET request
which inherits the original request properties (e.g. query string and headers).

```json
"Content": {
  "Type": "URL",
  "Request": [ "http://localhost:9000" ]
}
```

The complex form takes a JSON object, allowing custom request and response handling:

```json
"Content": {
  "Type": "URL",
  "Request": [{
    "URL": "http://localhost:9000",
    "Method": "GET",
    "Headers": {
      "$INHERIT": "1",
      "$REMOVE": [
        "X-Composure-Test"
      ],
      "$COPY": {
        "X-Composure-Test": "X-Composure-Copied"
      },
      "$SET": {
        "X-Composure": "1"
      },
      "$ADD": {
        "X-Foo": "1",
        "X-Bar": "2"
      }
    }
  }]
}
```

Header manipulation takes place in the following order:

* `$INHERIT` copies all request parameters
* `$REMOVE` removes a header
* `$COPY` copies a specific header (which can be renamed)
* `$SET` sets a header, overwriting it if it already exists
* `$ADD` adds a header, keeping existing headers with the same name

### Extending Composure

Composure is designed to be used as a standalone service.

It can also be extended using Go to provide additional functionality, for example:

* Inline pre- and post-flight transformations for faster response times
* Validating user authentication before invoking UI services
* Performing centralised rate limiting away from UI code

#### Load a composure

You can load a composure from an `io.File` object:

```go
f, _ := os.Open("spec.json")
spec, err := composition.Load(f)
```

Or you can parse JSON bytes directly:

```go
spec, err := composition.ParseJSON([]byte(`{...}`))
```

#### Render a composition

You can render a composition either with or without a HTTP request context.
The request context is important if you need to preserve or manipulate the
request headers or body.

To render without a context, simply pass in `nil`.

```go
func HandleRequest(w http.ResponseWriter, req *http.Request) {
  cm := spec.Get("/")
  err := cm.RenderFor(w, req)
}
```

For a more detailed example, see [main.go](main.go).

### Contributing

Clone this repository to ```$GOPATH/src/github.com/ian-kent/Composure```

Run `make` to install required dependencies and build Composure.

Run tests using ```make test``` or `goconvey`.

If you make any changes, run ```go fmt ./...``` before submitting a pull
request. This will be done automatically if you use `make`.

### Licence

Copyright ©‎ 2014, Ian Kent (http://www.iankent.eu).

Released under MIT license, see [LICENSE](license) for details.
