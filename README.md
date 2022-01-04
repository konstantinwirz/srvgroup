# srvgroup

[![Build Status](https://github.com/konstantinwirz/srvgroup/actions/workflows/main.yaml/badge.svg)](https://github.com/konstantinwirz/srvgroup/actions/workflows/main.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/konstantinwirz/srvgroup.svg)](https://pkg.go.dev/github.com/konstantinwirz/srvgroup)
[![Apache 2 licensed](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/oklog/run/master/LICENSE)

svcgroup.Group is a universal mechanism to manage the lifecycle of a server.

## Examples

### http.Server

```go
package main

import (
	"github.com/konstantinwirz/srvgroup"
	"log"
	"net/http"
)

func main() {
	appSrv := http.Server{Addr: ":8080"}
	metricsSrv := http.Server{Addr: ":8081"}

	errs := srvgroup.Run(
		srvgroup.HTTPServer(&appSrv),
		srvgroup.HTTPServer(&metricsSrv),
	)

	for _, err := range errs {
		log.Printf("error occurred: %v", err)
	}
}
```

### http.Server with lifecycle hooks

```go
package main

import (
	"github.com/konstantinwirz/srvgroup"
	"log"
	"net/http"
)

func main() {
	appSrv := http.Server{Addr: ":8080"}
	metricsSrv := http.Server{Addr: ":8081"}

	errs := srvgroup.Run(
		srvgroup.ServerLifecycleMiddleware(
			srvgroup.ServerLifecycleHooks{
				BeforeServe: func() {
					log.Printf("listening on port 8080...")
				},
				BeforeShutdown: func() {
					log.Printf("about to shutdown the server")
				},
			},
		)(srvgroup.HTTPServer(&appSrv)),
		srvgroup.HTTPServer(&metricsSrv),
	)

	for _, err := range errs {
		log.Printf("error occurred: %v", err)
	}
}
```

## Comparisons

Package srvgroup is heavily inspired by [oklog/run](https://github.com/oklog/run), it's almost the same implemnation,
the focus though is to manage lifecycle of servers which needs to be treated as a unit.