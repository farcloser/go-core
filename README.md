# Farcloser core go library

Provides basics to be used by all our go projects, specifically:
- logging (based on zerolog)
- reporter (based on Sentry)
- telemetry (based on otel)
- network (to ease manipulation of global network settings, specifically related to TLS)
- exec (to ease shelling out)
- human-readable units and duration
- filesystem safe operations primitives (locking, atomic writes)

## Dev

### Requirements

Review or call directly `./hack/dev-setup-linux.sh` or `./hack/dev-setup-macos.sh`

Then, `make install-linters`

### Makefile

```bash
make fix
make lint
make test
```

### Local documentation

```bash
go install golang.org/x/pkgsite/cmd/pkgsite@latest
pkgsite
open http://localhost:8080/go.farcloser.world/core
```

### Charter

1. This should contain solely stuff that is generic, likely to be used by 
any reasonnable go project (eg: logging, telemetry, etc.). Refrain from adding
here anything that is specific to a given project, or ecosystem.

2. Hide implementation away. Specifically, make sure the underlying dependencies
do not leak into your API.

### TODO

* consider going with https://github.com/go-logr/zerologr (eg: github.com/go-logr/logr)
if they really have traction - or slog alternatively
