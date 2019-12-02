
# Goroutine Report

`Goroutine Report` shows where and how much your Golang app creates goroutines.

```bash
    $ go get github.com/muravjov/goroutinereport
    $ curl http://localhost:6060/debug/pprof/goroutine?debug=2 | goroutinereport  | head -n 30
Created By                                                                                                                        Goroutine Count
                                                                                                                                  
(*Server).Serve:/usr/local/go/src/net/http/server.go:2927                                                                         363
(*Server).Serve:/usr/local/go/src/net/http/server.go:2927                                                                         92
(*Server).Serve:/usr/local/go/src/net/http/server.go:2927                                                                         28
(*Server).Serve:/usr/local/go/src/net/http/server.go:2927                                                                         10
(*connReader).startBackgroundRead:/usr/local/go/src/net/http/server.go:673                                                        9
(*Server).Serve:/usr/local/go/src/net/http/server.go:2927                                                                         2
init.0:/go/src/bitbucket.gcore.lu/con/concrete/vendor/k8s.io/klog/klog.go:404                                                     1
(*Server).Serve:/usr/local/go/src/net/http/server.go:2927                                                                         1
main:/go/src/bitbucket.gcore.lu/con/concrete/src/cementd/main.go:918                                                              1
::0                                                                                                                               1
main:/go/src/bitbucket.gcore.lu/con/concrete/src/cementd/main.go:924                                                              1
main:/go/src/bitbucket.gcore.lu/con/concrete/src/cementd/main.go:925                                                              1
OpenDB:/usr/local/go/src/database/sql/sql.go:722                                                                                  1
(*HTTPTransport).Configure.func1:/go/src/bitbucket.gcore.lu/con/concrete/vendor/github.com/getsentry/sentry-go/transport.go:152   1
OpenDB:/usr/local/go/src/database/sql/sql.go:723                                                                                  1

=======================================

363: select [0~11 minutes] [Created by http.(*Server).Serve @ server.go:2927]
    sql       /usr/local/go/src/database/sql/sql.go:1183                                                     (*DB).conn(#400, #25, #55, *, *, *, *)
    sql       /usr/local/go/src/database/sql/sql.go:1565                                                     (*DB).query(#400, #25, #55, #15, 0x2b, *, 0x1, 0x1, *, *, ...)
    sql       /usr/local/go/src/database/sql/sql.go:1547                                                     (*DB).QueryContext(#400, #25, #55, #15, 0x2b, *, 0x1, 0x1, 0, 0, ...)
    sql       /usr/local/go/src/database/sql/sql.go:1561                                                     (*DB).Query(...)
    state     /go/src/bitbucket.gcore.lu/con/concrete/src/state/advert.go:581                                checkNodeExists(*, #26, 0, 0)
    state     /go/src/bitbucket.gcore.lu/con/concrete/src/state/advert.go:19                                 DoNodeAdvert(*, #7, *)
    main      /go/src/bitbucket.gcore.lu/con/concrete/src/cementd/main.go:508                                (*serverHandler).ServeHTTP(#42, #2, *, *)
    http      /usr/local/go/src/net/http/server.go:2802                                                      serverHandler.ServeHTTP(#148, #2, *, *)
    http      /usr/local/go/src/net/http/server.go:1890                                                      (*conn).serve(*, #3, *)

```

## Syntax

```bash
$ goroutinereport [goroutines.dump]
```
