package source

import "io"
import "github.com/gozips/source"
import "github.com/gozips/sources"

// LimitedCloser is a wrapper around LimitedReader to implement ReadCloser
type LimitedCloser struct {
	*io.LimitedReader
	io.ReadCloser
}

func (l LimitedCloser) Read(b []byte) (int, error) {
	return l.LimitedReader.Read(b)
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
