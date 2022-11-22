# GAPI Addons

A repo with a list of additional packages to improve GAPI features :

| Description                | Doc                                                                                                                                                             |
|----------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Auth middleware for GCloud | [![PkgGoDev](https://pkg.go.dev/badge/github.com/mwm-io/gapi-addons/middleware/auth)](https://pkg.go.dev/github.com/mwm-io/gapi-addons/middleware/auth) |


## How it works

First choose your package(s) and add it using go modules, here is an example for [auth.GCloudServiceAccount](https://pkg.go.dev/github.com/mwm-io/gapi-addons/gcloud/middleware/auth#GCloudServiceAccount): 
```sh
$ go get github.com/mwm-io/gapi-addons/gcloud/middleware/auth@latest
```

The use it in you code:

```go
package main

import (
    ...
    github.com/mwm-io/gapi-addons/gcloud/middleware/auth
)

type myHandler struct {
    handler.WithMiddlewares
}

// NewHandler is your handler factory
func NewHandler() handler.Handler {
    var h myHandler
	
    h.MiddlewareList = []handler.Middleware{
        auth.GCloudServiceAccount{
            ServiceAccount: "lorem-ispum@gserviceaccount.google.com",
        },
    }
    ...
}
```

