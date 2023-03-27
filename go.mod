module github.com/gastrodon/popplio

go 1.19

require (
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/jcelliott/turnpike v0.0.0-20210629143239-1dadcad507a3
)

require github.com/mitchellh/mapstructure v1.5.0

require (
	github.com/google/uuid v1.3.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	golang.org/x/sys v0.6.0 // indirect
)

replace github.com/jcelliott/turnpike => ../turnpike
