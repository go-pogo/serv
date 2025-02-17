serv
====
[![Latest release][latest-release-img]][latest-release-url]
[![Build status][build-status-img]][build-status-url]
[![Go Report Card][report-img]][report-url]
[![Documentation][doc-img]][doc-url]

[latest-release-img]: https://img.shields.io/github/release/go-pogo/serv.svg?label=latest

[latest-release-url]: https://github.com/go-pogo/serv/releases

[build-status-img]: https://github.com/go-pogo/serv/actions/workflows/test.yml/badge.svg

[build-status-url]: https://github.com/go-pogo/serv/actions/workflows/test.yml

[report-img]: https://goreportcard.com/badge/github.com/go-pogo/serv

[report-url]: https://goreportcard.com/report/github.com/go-pogo/serv

[doc-img]: https://godoc.org/github.com/go-pogo/serv?status.svg

[doc-url]: https://pkg.go.dev/github.com/go-pogo/serv


Package `serv` provides a server and router implementation based on the `http`
package, with a focus on security, flexibility and ease of use.

Included features:
- `Server` with sane and safe defaults;
- `Server` `State` retrieval;
- `Router`/`ServeMux` with easy (mass) `Route` registration;
- Set custom "not found" `http.Handler` on `ServeMux`;
- support for access logging.

<hr>

```sh
go get github.com/go-pogo/serv
```

```go
import "github.com/go-pogo/serv"
```

## Documentation

Additional detailed documentation is available at [pkg.go.dev][doc-url]

## Created with

<a href="https://www.jetbrains.com/?from=go-pogo" target="_blank"><img src="https://resources.jetbrains.com/storage/products/company/brand/logos/GoLand_icon.png" width="35" /></a>

## License

Copyright © 2021-2025 [Roel Schut](https://roelschut.nl). All rights reserved.

This project is governed by a BSD-style license that can be found in the [LICENSE](LICENSE) file.
