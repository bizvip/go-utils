package obs

import (
	"errors"
	"io"
)

type ICallbackReadCloser interface {
	setCallbackReadCloser(body io.ReadCloser)
}

func (output *PutObjectOutput) setCallbackReadCloser(body io.ReadCloser) {
	output.CallbackBody.data = body
}

func (output *CompleteMultipartUploadOutput) setCallbackReadCloser(body io.ReadCloser) {
	output.CallbackBody.data = body
}

// define CallbackBody
type CallbackBody struct {
	data io.ReadCloser
}

func (output CallbackBody) ReadCallbackBody(p []byte) (int, error) {
	if output.data == nil {
		return 0, errors.New("have no callback data")
	}
	return output.data.Read(p)
}

func (output CallbackBody) CloseCallbackBody() error {
	if output.data == nil {
		return errors.New("have no callback data")
	}
	return output.data.Close()
}
