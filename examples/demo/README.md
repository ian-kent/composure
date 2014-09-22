Composure demo
==============

- Run `composure` to start the Composure demo
- Open it in a browser: http://localhost:9000
- Run it in [Docker](https://www.docker.com/) with the [example Dockerfile](Dockerfile)

### Overview

The Composure demo uses three Gists to build the page.

- It loads and renders the Navbar and Content components
- It loads and renders the Template component, injecting Navbar and Content
- It adds a custom Server header in the Postflight phase
