# gozips-sources

[![Build Status](https://travis-ci.org/nowk/gozips-sources.svg?branch=master)](https://travis-ci.org/nowk/gozips-sources)
[![GoDoc](https://godoc.org/github.com/nowk/gozips-sources?status.svg)](http://godoc.org/github.com/nowk/gozips-sources)

Source funcs for [gozips](https://github.com/gozips)

## HTTPLimit

    out := new(bytes.Buffer)
    zip := gozips.NewZip(sources.HTTPLimit(1024))
    zip.Add(url1)
    zip.Add(url2, url3)
    n, err := zip.WriteTo(out)

    // each entry written will be truncated to 1024 bytes

## License

MIT