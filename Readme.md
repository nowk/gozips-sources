# gozip-sources

[![Build Status](https://travis-ci.org/nowk/gozip-sources.svg?branch=master)](https://travis-ci.org/nowk/gozip-sources)
[![GoDoc](https://godoc.org/github.com/nowk/gozip-sources?status.svg)](http://godoc.org/github.com/nowk/gozip-sources)

Source funcs for [gozip](https://github.com/gozips)

## HTTPlimited

    out := new(bytes.Buffer)
    zip := NewZip(source.HTTPlimited(1024))
    zip.Add(url1)
    zip.Add(url2, url3)
    n, err := zip.WriteTo(out)

    // each entries written will be truncated to 1024 bytes

## License

MIT