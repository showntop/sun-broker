package server

import (
	"errors"
	"fmt"
)

var UnsupportedProtocol = errors.New("unsupported protocol")
var UnparsedServerUrl = errors.New("can not parse the server url")

type EorroCode uint8

const (
	_ EorroCode = iota
	LaunchError
	ListentError
	AcceptError
)

var errorMessageMap = map[EorroCode]string{
	LaunchError:  "server start error",
	ListentError: "network error when listener",
	AcceptError:  "accept conn error",
}

type serverError struct {
	err  error
	code EorroCode
}

func (se *serverError) Error() string {
	return fmt.Sprintf("%s: %s", errorMessageMap[se.code], se.err.Error())
}

func newServerError(code EorroCode, err error) *serverError {
	return &serverError{
		code: code,
		err:  err,
	}
}
