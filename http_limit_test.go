package source

import "bytes"
import "fmt"
import "io"
import "net/http"
import "net/http/httptest"
import "regexp"
import "testing"
import "github.com/nowk/assert"

func h(str string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(str))
	}
}

func tServer() (ts *httptest.Server) {
	mux := http.NewServeMux()
	mux.HandleFunc("/index.html", h("Hello World!"))
	mux.HandleFunc("/posts", h("Post Body"))
	mux.HandleFunc("/api/data.json", h(`{"data": ["one"]}`))
	ts = httptest.NewServer(mux)

	return
}

func TestFileIsWithinLimits(t *testing.T) {
	ts := tServer()
	defer ts.Close()

	for _, v := range []struct {
		u, n, m string
	}{
		{"/index.html", "index.html", "Hello World!"},
		{"/posts", "posts", "Post Body"},
		{"/api/data.json", "data.json", `{"data": ["one"]}`},
	} {
		url := fmt.Sprintf("%s%s", ts.URL, v.u)
		l := int64(len(v.m))
		name, r, _ := HTTPLimit(l)(url)
		defer r.Close()

		var b []byte
		buf := bytes.NewBuffer(b)
		n, err := io.Copy(buf, r)

		assert.Nil(t, err)
		assert.Equal(t, l, n)
		assert.Equal(t, v.n, name)
		assert.Equal(t, v.m, buf.String())
	}
}

func TestFileExeedsByteLimit(t *testing.T) {
	ts := tServer()
	defer ts.Close()

	for _, v := range []struct {
		u, n, m string
	}{
		{"/index.html", "index.html", "Hello World!"},
		{"/posts", "posts", "Post Body"},
		{"/api/data.json", "data.json", `{"data": ["one"]}`},
	} {
		url := fmt.Sprintf("%s%s", ts.URL, v.u)
		l := int64(len(v.m)) - 2
		name, r, _ := HTTPLimit(l)(url)
		defer r.Close()

		var b []byte
		buf := bytes.NewBuffer(b)
		n, err := io.Copy(buf, r)

		assert.Equal(t, "error: limit: exceeded allowable read limit", err.Error())
		assert.TypeOf(t, "source.ReadError", err)
		assert.Equal(t, l, n)
		assert.Equal(t, v.n, name)
		assert.Equal(t, v.m[:l], buf.String())
	}
}

func TestHTTPClientError(t *testing.T) {
	reg := regexp.MustCompile(`Get http:\/\/unreachable:( dial tcp:)? lookup unreachable: no such host`)

	// fails if ISP picks up and redirects to search, which TWC does
	name, v, err := HTTPLimit(4)("http://unreachable")
	assert.Equal(t, "unreachable.txt", name)
	if !reg.MatchString(err.Error()) {
		t.Errorf("Expected %s to match %s", err.Error(), reg.String())
	}

	b := make([]byte, 32*1024)
	r := v.(io.ReadCloser)
	n, _ := r.Read(b)
	if str := string(b[:n]); !reg.MatchString(str) {
		t.Errorf("Expected %s to match %s", str, reg.String())
	}

	r.Close()
}
