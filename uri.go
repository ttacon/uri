package uri

import (
	"bytes"
	"errors"
	"strings"
)

// URI represents a general RFC3986 specified URI.
type URI interface {
	// Scheme is the scheme the URI conforms to.
	Scheme() string

	// Authority returns the authority information for the URI.
	Authority() Authority

	// QueryPieces returns a map of key/value pairs of all parameters
	// in the query string of the URI.
	QueryPieces() map[string]string

	// Fragment returns the fragement (component proceeded by '#') in the
	// URI if there is one.
	Fragment() string

	// Builder returns a Builder that can be used to modify the URI.
	Builder() URIBuilder

	String() string
}

// Authority represents the authority information that a URI can contained
// as specified by RFC3986. Be aware that the RFC does not specify any
// information on username/password formatting - both are contained within
// UserInfo.
type Authority interface {
	UserInfo() string
	Host() string
	Port() string
	String() string
}

var (
	colon    = ":"
	question = "?"
	fragment = "#"
)

// ParseURI attempts to parse a URI and only returns an error if the URI
// is not RFC3986 compliant.
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
	prefix   string
	userinfo string
	host     string
	port     string
	path     string
}

func (a authorityInfo) UserInfo() string { return a.userinfo }
func (a authorityInfo) Host() string     { return a.host }
func (a authorityInfo) Port() string     { return a.port }
func (a authorityInfo) String() string {
	var buf = bytes.NewBuffer(nil)
	buf.WriteString(a.prefix)
	buf.WriteString(a.userinfo)
	if len(a.userinfo) > 0 {
		buf.Write(atBytes)
	}
	buf.WriteString(a.host)
	if len(a.port) > 0 {
		buf.Write(colonBytes)
	}
	buf.WriteString(a.port)
	buf.WriteString(a.path)
	return buf.String()
}

var (
	// byte literals
	atBytes       = []byte("@")
	colonBytes    = []byte(":")
	fragmentBytes = []byte("#")
	queryBytes    = []byte("?")
	slashBytes    = []byte("/")
)

func parseAuthority(hier string) *authorityInfo {
	// as per RFC 3986 Section 3.6
	var prefix, userinfo, host, port, path string

	// authority sections can begin with a '//'
	if strings.HasPrefix(hier, "//") {
		prefix = "//"
		hier = strings.TrimPrefix(hier, "//")
	}

	slashEnd := strings.Index(hier, "/")
	if slashEnd > 0 {
		if slashEnd < len(hier) {
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
		prefix:   prefix,
		userinfo: userinfo,
		host:     host,
		port:     port,
		path:     path,
	}
}

// Builder time!

// URIBuilder is a construct for building URIs.
type URIBuilder interface {
	SetScheme(scheme string) URIBuilder
	SetUserInfo(userinfo string) URIBuilder
	SetHost(host string) URIBuilder
	SetPort(port string) URIBuilder
	SetPath(path string) URIBuilder
	SetQuery(query string) URIBuilder
	SetFragment(fragment string) URIBuilder

	// Returns the URI this URIBuilder represents.
	Build() URI
	String() string
}

func (u *uri) ensureAuthorityExists() {
	if u.authority == nil {
		u.authority = &authorityInfo{}
	}
}

func (u *uri) SetScheme(scheme string) URIBuilder {
	u.scheme = scheme
	return u
}

func (u *uri) SetUserInfo(userinfo string) URIBuilder {
	u.ensureAuthorityExists()
	u.authority.userinfo = userinfo
	return u
}

func (u *uri) SetHost(host string) URIBuilder {
	u.ensureAuthorityExists()
	u.authority.host = host
	return u
}

func (u *uri) SetPort(port string) URIBuilder {
	u.ensureAuthorityExists()
	u.authority.port = port
	return u
}

func (u *uri) SetPath(path string) URIBuilder {
	u.ensureAuthorityExists()
	u.authority.path = path
	return u
}

func (u *uri) SetQuery(query string) URIBuilder {
	u.query = query
	return u
}

func (u *uri) SetFragment(fragment string) URIBuilder {
	u.fragment = fragment
	return u
}

func (u *uri) Build() URI {
	return u
}

func (u *uri) Builder() URIBuilder {
	return u
}

func (u *uri) String() string {
	// TODO(ttacon): how do we know if // exists in authority info?
	// keep prefix on authority info?
	var buf = bytes.NewBuffer(nil)
	if len(u.scheme) > 0 {
		buf.WriteString(u.scheme)
		buf.Write(colonBytes)
	}

	buf.WriteString(u.authority.String())

	if len(u.query) > 0 {
		buf.Write(queryBytes)
		buf.WriteString(u.query)
	}

	if len(u.fragment) > 0 {
		buf.Write(fragmentBytes)
		buf.WriteString(u.fragment)
	}

	return buf.String()
}
