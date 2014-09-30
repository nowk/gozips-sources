# gozips-sources

[![Build Status](https://travis-ci.org/nowk/gozips-sources.svg?branch=master)](https://travis-ci.org/nowk/gozips-sources)
[![GoDoc](https://godoc.org/github.com/nowk/gozips-sources?status.svg)](http://godoc.org/github.com/nowk/gozips-sources)

Source funcs for [gozips](https://github.com/gozips)

## HTTPlimited

    out := new(bytes.Buffer)
    zip := NewZip(source.HTTPlimited(1024))
    zip.Add(url1)
    zip.Add(url2, url3)
    n, err := zip.WriteTo(out)

    // each entries written will be truncated to 1024 bytes

## License

MIT