# uri
[![Build Status](https://travis-ci.org/fredbi/uri.svg?branch=master)](https://travis-ci.org/fredbi/uri)
[![codecov](https://codecov.io/gh/fredbi/uri/branch/master/graph/badge.svg)](https://codecov.io/gh/fredbi/uri)
[![license](http://img.shields.io/github/license/mapshape/apistatus.svg)](https://raw.githubusercontent.com/fredbi/uri/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/go-openapi/strfmt?status.svg)](http://godoc.org/github.com/go-openapi/strfmt)
[![GolangCI](https://golangci.com/badges/github.com/fredbi/uri.svg)](https://golangci.com)
[![Go Report Card](https://goreportcard.com/badge/github.com/fredbi/uri)](https://goreportcard.com/report/github.com/fredbi/uri)

Package uri is meant to be an RFC 3986 compliant URI builder, parser and validator for golang.

References:
* https://tools.ietf.org/html/rfc3986

Internationalization support:
* https://tools.ietf.org/html/rfc3987

IPv6 addressing scheme reference and erratum:
* https://tools.ietf.org/html/rfc6874

This allows for stricter conformance than the `net/url` golang standard libary,
which provides a workable but loose implementation of the RFC.

This package concentrates on RFC 3986 strictness for URI validation. 
At the moment, there is no attempt to normalize or auto-escape strings. 
For url normalization, see github.com/PuertokitoBio/purell.

## Disclaimer

Not supported:
* provisions for "IPvFuture" are not implemented

## Credits

Tests have been aggregated from test suites of URI validators from other languages:
perl, python, scala, .Net.

> This package is based on the work from ttacon/uri (credits: Trey Tacon).
> Extra features like MySQL URIs present in the original repo have been removed.
