# CodeComet core go library

Provides basics to be used by all our go projects, specifically:
- logging (based on zerolog)
- reporter (based on Sentry)
- network (to ease manipulation of global network settings, specifically related to TLS)
- exec (to ease shelling out)

## Dev

### Makefile

* make lint
* make lint-fix
* make tidy

### Local documentation

Install pkgsite: go install golang.org/x/pkgsite/cmd/pkgsite@latest

Run it: pkgsite

Open: http://localhost:8080/github.com/go.codecomet.dev/core
