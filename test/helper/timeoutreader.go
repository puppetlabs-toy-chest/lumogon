package helper

import (
	"errors"
	"io"
)

// ErrTimeout is thrown when the number of valid reads is exceeded
var ErrTimeout = errors.New("timeout")

// TimeoutReader wraps a Reader and allows up to n valid reads after
// which subsequent reads will return ErrTimeout
func TimeoutReader(r io.Reader, n int) io.Reader { return &timeoutReader{r, 0, n} }

type timeoutReader struct {
	r          io.Reader
	count      int
	validReads int
}

func (r *timeoutReader) Read(p []byte) (int, error) {
	r.count++
	if r.count > r.validReads {
		return 0, ErrTimeout
	}
	return r.r.Read(p)
}
