To-do
=====

- HTTP component URL placeholders, e.g. to proxy an entire node
- Handle query parameters and request body
- Load template asynchronously, just defer rendering
- Clean up documentation (particularly components and URL request/responses)
- Consider lazy evaluation of components
- Performance
  - Compile JSON specs to go functions
- Component XPath support
  - Selectively include parts of a component
  - Reuse the component to select different parts
- HTTP caching
  - Something like https://github.com/gregjones/httpcache
  - Respect HTTP headers so components control caching
  - Consider implication of caching rendered compositions
- Golang integration demo (e.g. parsing Authorization header)
- More spec demos (preflight, postflight, etc)
- HTTP API for spec manipulation/reloading
- Host matching for routes
- Change default HTTP component behaviour
  - set $INHERIT for query and headers
  - placeholders into URLs
