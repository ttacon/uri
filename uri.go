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
	var (
		schemeEnd   = strings.Index(raw, colon)
		hierPartEnd = strings.Index(raw, question)
		queryEnd    = strings.Index(raw, fragment)

		curr int
	)

	if schemeEnd < 1 {
		return nil, ErrNoSchemeFound
	} else if schemeEnd+1 == len(raw) {
		return nil, ErrInvalidURI
	}

	scheme := raw[curr:schemeEnd]
	curr = schemeEnd + 1

	if hierPartEnd == len(raw) || (hierPartEnd < 0 && queryEnd < 0) {
		return &uri{
			scheme:    scheme,
			hierPart:  raw[curr:], // since it's all that's left
			authority: parseAuthority(raw[curr:]),
		}, nil
	}

	var (
		hierPart, query, fragment string
		authorityInfo             *authorityInfo
	)
	if hierPartEnd > 0 {
		hierPart = raw[curr:hierPartEnd]
		authorityInfo = parseAuthority(hierPart)
		if hierPartEnd+1 < len(raw) {
			query = raw[hierPartEnd+1:]
		}
		curr = hierPartEnd + 1
	}

	if queryEnd == len(raw) {
		if hierPartEnd < 0 {
			hierPart = raw[curr:queryEnd]
			authorityInfo = parseAuthority(hierPart)
		} else {
			query = raw[curr:hierPartEnd]
			hierPart = raw[hierPartEnd+1 : queryEnd]
			authorityInfo = parseAuthority(hierPart)
		}
		return &uri{
			scheme:    scheme,
			hierPart:  hierPart,
			authority: authorityInfo,
			query:     query,
		}, nil
	} else if queryEnd > 0 {
		if hierPartEnd < 0 {
			hierPart = raw[curr:queryEnd]
			authorityInfo = parseAuthority(hierPart)
		} else {
			query = raw[curr:queryEnd]
		}
		if queryEnd+1 < len(raw) {
			fragment = raw[queryEnd+1:]
		}
	}

	return &uri{
		scheme:    scheme,
		hierPart:  hierPart,
		query:     query,
		fragment:  fragment,
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
	authority *authorityInfo
}

func (u *uri) Scheme() string {
	return u.scheme
}

func (u *uri) Authority() Authority {
	return u.authority
}

func (u *uri) QueryPieces() map[string]string {
	return nil
}

func (u *uri) Fragment() string {
	return u.fragment
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

func (a authorityInfo) UserInfo() string { return a.userinfo }
func (a authorityInfo) Host() string     { return a.host }
func (a authorityInfo) Port() string     { return a.port }

func parseAuthority(hier string) *authorityInfo {
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

	return &authorityInfo{
		userinfo: userinfo,
		host:     host,
		port:     port,
		path:     path,
	}
}

// Builder time!

type URIBuilder interface {
	SetScheme(scheme string)
	SetUserInfo(userinfo string)
	SetHost(host string)
	SetPort(port string)
	SetPath(path string)
	SetQuery(query string)
	SetFragment(fragment string)

	Build() string
	String() string
}
