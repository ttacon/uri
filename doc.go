// Package uri is meant to be an RFC 3986 compliant URI builder and parser.
//
// I built this package because I was tired of having to always
// build conditional URI building for connection to databses (for
// servers and tools). The two main components are the uri parser
// which allows you to parse a URI and extract information from it and
// the builder which allows you to build URIs in a simple, programmatic
// fashion. I also added two convenienve builders for the gomysql and postgres
// drivers. Feel free to add more as desired!
package uri
