package uri

import (
	"errors"
	"strings"
)

type URI interface {
	Scheme() string
	Authority() Authority

	QueryPieces() map[string]string
	Fragment() string
}

type Authority interface {
	UserInfo() string
	Host() string
	Port() string
}

var (
	colon    = ":"
	question = "?"
	fragment = "#"
)

func ParseURI(raw string) (URI, error) {
	schemeEnd := strings.Index(raw, colon)
	if schemeEnd < 1 {
		return nil, ErrNoSchemeFound
	} else if schemeEnd == len(raw) {
		return nil, ErrInvalidURI
	}

	scheme := raw[:schemeEnd]
	raw = raw[schemeEnd+1:]

	hierPartEnd := strings.Index(raw, question)
	if hierPartEnd < 0 || hierPartEnd == len(raw) {
		return &uri{
			scheme:   scheme,
			hierPart: raw, // since it's all that's left
		}, nil
	}

	hierPart := raw[:hierPartEnd]
	raw = raw[hierPartEnd+1:]

	queryEnd := strings.Index(raw, fragment)
	if queryEnd < 0 || queryEnd == len(raw) {
		return &uri{
			scheme:   scheme,
			hierPart: hierPart,
			query:    raw,
		}, nil
	}

	return &uri{
		scheme:   scheme,
		hierPart: hierPart,
		query:    raw[:queryEnd],
		fragment: raw[queryEnd+1:],
	}, nil
}

type uri struct {
	// raw components
	scheme   string
	hierPart string
	query    string
	fragment string
}

func (u *uri) Scheme() string {
	return ""
}

func (u *uri) Authority() Authority {
	return nil
}

func (u *uri) QueryPieces() map[string]string {
	return nil
}

func (u *uri) Fragment() string {
	return ""
}

// errors
var (
	ErrNoSchemeFound = errors.New("no scheme found in URI")
	ErrInvalidURI    = errors.New("not a valid URI")
)
