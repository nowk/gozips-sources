package source

import "bytes"
import "fmt"
import "io"
import "net/http"
import "net/http/httptest"
import "regexp"
import "testing"
import "github.com/gozips/zips"
import gozipt "github.com/gozips/testing"
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

func TestLimitsBodyToN(t *testing.T) {
	ts := tServer()
	defer ts.Close()

	url1 := fmt.Sprintf("%s/index.html", ts.URL)
	url2 := fmt.Sprintf("%s/posts", ts.URL)
	url3 := fmt.Sprintf("%s/api/data.json", ts.URL)

	out := new(bytes.Buffer)
	zip := zips.NewZip(HTTPLimit(2))
	zip.Add(url1)
	zip.Add(url2, url3)
	n, err := zip.WriteTo(out)

	assert.Nil(t, err)
	assert.Equal(t, int64(6), n)
	gozipt.VerifyZip(t, out.Bytes(), []gozipt.Entries{
		{"index.html", "He"},
		{"posts", "Po"},
		{"data.json", `{"`},
	})
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
