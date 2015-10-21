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
			scheme:    scheme,
			hierPart:  raw, // since it's all that's left
			authority: parseAuthority(raw),
		}, nil
	}

	hierPart := raw[:hierPartEnd]

	// parse authority
	authorityInfo := parseAuthority(hierPart)

	raw = raw[hierPartEnd+1:]

	queryEnd := strings.Index(raw, fragment)
	if queryEnd < 0 || queryEnd == len(raw) {
		return &uri{
			scheme:    scheme,
			hierPart:  hierPart,
			query:     raw,
			authority: authorityInfo,
		}, nil
	}

	return &uri{
		scheme:    scheme,
		hierPart:  hierPart,
		query:     raw[:queryEnd],
		fragment:  raw[queryEnd+1:],
		authority: authorityInfo,
	}, nil
}

type uri struct {
	// raw components
	scheme   string
	hierPart string
	query    string
	fragment string

	// parsed components
	authority authorityInfo
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

type authorityInfo struct {
	userinfo string
	host     string
	port     string
	path     string
}

func parseAuthority(hier string) authorityInfo {
	// as per RFC 3986 Section 3.6
	var userinfo, host, port, path string

	// authority sections can begin with a '//'
	hier = strings.TrimPrefix(hier, "//")

	slashEnd := strings.Index(hier, "/")
	if slashEnd > 0 {
		if slashEnd+1 < len(hier) {
			path = hier[slashEnd:]
		}
		hier = hier[:slashEnd]
	}

	// authority   = [ userinfo "@" ] host [ ":" port ]
	host = hier
	if at := strings.Index(host, "@"); at > 0 {
		userinfo = host[:at]

		if at+1 < len(host) {
			host = host[at+1:]
		}
	}

	if colon := strings.Index(host, ":"); colon > 0 {
		if colon+1 < len(host) {
			port = host[colon+1:]
		}

		host = host[:colon]
	}

	return authorityInfo{
		userinfo: userinfo,
		host:     host,
		port:     port,
		path:     path,
	}
}
