package source

import "io"
import "github.com/gozips/source"
import "github.com/gozips/sources"

// LimitedCloser is a wrapper around LimitedReader to implement ReadCloser
type LimitedCloser struct {
	*io.LimitedReader
	io.ReadCloser
}

// Read delegates to LimitedReader, at EOF it tries to read an extra byte. If
// the read returns 1, means file exceeds the limit and source.Error is
// returned
func (l LimitedCloser) Read(b []byte) (int, error) {
	n, err := l.LimitedReader.Read(b)
	if err == io.EOF {
		r := l.LimitedReader.R
		m, _ := r.Read(make([]byte, 1))
		if m > 0 {
			return n, source.ReadError{
				Message: "error: limit: exceeded allowable read limit",
			}
		}
	}

	return n, err
}

// HTTPLimit returns an http body that reads only up to n
func HTTPLimit(n int64) source.Func {
	return func(urlStr string) (string, io.ReadCloser, error) {
		name, r, err := sources.HTTP(urlStr)
		if err != nil {
			return name, r, err
		}

		l := &io.LimitedReader{r, n}
		c := LimitedCloser{
			l,
			r,
		}

		return name, c, err
	}
}
