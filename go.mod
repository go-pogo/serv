module github.com/go-pogo/serv

go 1.23.0

// cors subpackage was not working properly, do not use it
retract v0.6.0

require (
	github.com/felixge/httpsnoop v1.0.4
	github.com/go-pogo/easytls v0.1.4
	github.com/go-pogo/errors v0.12.0
	github.com/go-pogo/rawconv v0.6.4
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
