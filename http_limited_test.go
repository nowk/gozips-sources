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
	zip := zips.NewZip(HTTPlimited(2))
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
	errmsg := "Get http://unreachable: dial tcp: lookup unreachable: no such host"

	// fails if ISP picks up and redirects to search, which TWC does
	name, v, err := HTTPlimited(4)("http://unreachable")
	assert.Equal(t, "unreachable.txt", name)
	assert.Equal(t, errmsg, err.Error())

	b := make([]byte, 32*1024)
	r := v.(io.ReadCloser)
	n, _ := r.Read(b)
	if str := string(b[:n]); !reg.MatchString(str) {
		t.Errorf("Expected %s, got %s", reg.String(), str)
	}

	r.Close()
}
