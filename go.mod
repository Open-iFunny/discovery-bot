module github.com/gastrodon/popplio

go 1.19

require (
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/jcelliott/turnpike v0.0.0-20210629143239-1dadcad507a3
)

require github.com/mitchellh/mapstructure v1.5.0

require github.com/ugorji/go/codec v1.2.11 // indirect

replace github.com/jcelliott/turnpike => ../turnpike
