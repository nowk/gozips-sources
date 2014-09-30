package source

import "io"
import "github.com/gozips/source"

// LimitedCloser is a wrapper around LimitedReader to implement ReadCloser
type LimitedCloser struct {
	*io.LimitedReader
	io.ReadCloser
}

func (l LimitedCloser) Read(b []byte) (int, error) {
	return l.LimitedReader.Read(b)
}

// HTTPlimited returns an http body that reads only up to n
func HTTPlimited(n int64) func(string) (string, interface{}) {
	return func(urlStr string) (string, interface{}) {
		name, r := source.HTTP(urlStr)

		switch v := r.(type) {
		case io.ReadCloser:
			l := &io.LimitedReader{v, n}
			c := LimitedCloser{
				l,
				v,
			}
			return name, c

		case error:
			return name, v
		}

		return "", nil
	}
}
