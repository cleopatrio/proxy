package core

import (
	"bytes"
	"io"
)

type RequestBody struct {
	Data []byte

	state bytes.Buffer
}

func (rb *RequestBody) Close() error {
	return io.NopCloser(&rb.state).Close()
}

func (rb *RequestBody) Read(buffer []byte) (int, error) {
	return rb.state.Read(buffer)
}

var _ io.ReadCloser = (*RequestBody)(nil)
