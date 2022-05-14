# Uber Fx Helpers

Some helper for the https://github.com/uber-go/fx

[![CI](https://github.com/albenik/fx-glue/actions/workflows/ci.yaml/badge.svg)](https://github.com/albenik/fx-glue/actions/workflows/ci.yaml)

## Install

```shell
go get -u github.com/albenik/fx-glue
```

## Use

```go
package main

import (
    "net/http"

    "github.com/albenik/fx-glue"
    "github.com/albenik/huenv"
    "github.com/rs/zerolog"
    "github.com/uber-go/fx"
)

type Config struct {
    App        fxglue.AppConfig
    ListenAddr string `fx:"supply,name=listen_addr"`
}

func main() {
    conf := new(Config)
    if err := huenv.Init(conf); err != nil {
		panic(err)
    }

    mod := fx.Module("app",
        fxglue.SupplyConfig(conf),
        fx.Provide(
            func() http.Handler {
                return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                    // ...
                })
            },
        ),
        fx.Invoke(listenAndServer),
    )

    app := fx.New(mod)
    app.Run()
}

type listenAndServeArgs struct {
    fx.In

    ListenAddr string `name:"listen_addr"`
    Handler    http.Handler
    Log        *zerolog.Logger
}

func listenAndServe(lc fx.Lifecycle, args listenAndServeArgs) {
    srv := &http.Server{
        Addr:    args.ListenAddr,
        Handler: args.Handler,
    }
    lc.Append(fxglue.NewHTTPServerHook(srv, args.Log))
}
```
