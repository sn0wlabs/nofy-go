package client

import "io"

//
// Private types (stolen from stripe-go)
//

// nopReadCloser's sole purpose is to give us a way to turn an `io.Reader` into
// an `io.ReadCloser` by adding a no-op implementation of the `Closer`
// interface. We need this because `http.Request`'s `Body` takes an
// `io.ReadCloser` instead of a `io.Reader`.
type nopReadCloser struct {
	io.Reader
}

func (nopReadCloser) Close() error { return nil }

//
// Response Types
//

type Res[T any] struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Data    *T     `json:"data"`
}

//
// Error Type
//

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var (
	ParamEncodeError = Error{
		Code:    "paramEncodeError",
		Message: "There was an error encoding your parameters.",
	}
	RequestError = Error{
		Code:    "requestError",
		Message: "There was an error with your request.",
	}
	MultipartRequestError = Error{
		Code:    "multipartRequestError",
		Message: "There was an error with your multipart request.",
	}
	NoResponseTypeError = Error{
		Code:    "noResponseTypeError",
		Message: "There was no response type specified.",
	}
)

type Map map[string]interface{}
