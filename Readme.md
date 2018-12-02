Package alb-gateway is based on a fork of [Apex/gateway](https://github.com/apex/gateway)

This package provides a drop-in replacement for net/http's `ListenAndServe` for use in an AWS Lambda invoked by a ALB 
Labmda Target Group, simply swap it out for `gateway.ListenAndServe`.

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/getas/alb-gateway"
	"github.com/aws/aws-lambda-go"
)

func main() {
	http.HandleFunc("/", hello)
	log.Fatal(gateway.ListenAndServe(":3000", nil))
}

func hello(w http.ResponseWriter, r *http.Request) {
	// example retrieving values from the api gateway proxy request context.
	_, ok := gateway.RequestContext(r.Context())
	if !ok {
		fmt.Fprint(w, "Hello World from Go")
		return
	}

	fmt.Fprintf(w, "Hello %s from Go", r.RemoteAddr)
}
```

---

![](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/status-stable-green.svg)
