// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

/*
Package serv provides a server and router implementation based on the http
package, with a focus on security, flexibility and ease of use.

Included features:
- [Server] with sane and safe defaults;
- [Server] [State] retrieval;
- [Router]/[ServeMux] with easy (mass) [Route] registration;
- Set custom "not found" http.Handler on [ServeMux];
- support for access logging;
- default TLS configuration.
*/
