package uri

import (
	"bytes"
	"errors"
	"net"
	"regexp"
	"strings"
)

// errors
var (
	ErrNoSchemeFound    = errors.New("no scheme found in URI")
	ErrInvalidURI       = errors.New("not a valid URI")
	ErrInvalidCharacter = errors.New("invalid character in URI")
	ErrInvalidScheme    = errors.New("invalid scheme in URI")
	ErrInvalidQuery     = errors.New("invalid query string in URI")
	ErrInvalidFragment  = errors.New("invalid fragment in URI")
	ErrInvalidPath      = errors.New("invalid path in URI")
	ErrInvalidHost      = errors.New("invalid host in URI")
	ErrInvalidPort      = errors.New("invalid port in URI")
	ErrInvalidUserInfo  = errors.New("invalid userinfo in URI")
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
	Builder() Builder

	// String return a string representation of the URI
	String() string

	// Validate the different components of the URI
	Validate() error
}

// Authority represents the authority information that a URI may contain
// as specified by RFC3986. Be aware that the RFC does not specify any
// information on username/password formatting - both are contained within
// UserInfo.
type Authority interface {
	UserInfo() string
	Host() string
	Port() string
	Path() string
	String() string
	Validate(...string) error
}

// Builder is a construct for building URIs.
type Builder interface {
	SetScheme(scheme string) Builder
	SetUserInfo(userinfo string) Builder
	SetHost(host string) Builder
	SetPort(port string) Builder
	SetPath(path string) Builder
	SetQuery(query string) Builder
	SetFragment(fragment string) Builder

	// Returns the URI this Builder represents.
	Build() URI
	String() string
}

const (
	colon    = ":"
	question = "?"
	fragment = "#"
)

// IsURI tells if a URI is valid according to RFC3986
func IsURI(raw string) bool {
	_, err := ParseURI(raw)
	return err == nil
}

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
		authorityInfo, err := parseAuthority(raw[curr:])
		if err != nil {
			return nil, ErrInvalidURI
		}
		u := &uri{
			scheme:    scheme,
			hierPart:  raw[curr:], // since it's all that's left
			authority: authorityInfo,
		}
		return u, u.Validate()
	}

	var (
		hierPart, query, fragment string
		authorityInfo             *authorityInfo
	)
	var err error
	if hierPartEnd > 0 {
		hierPart = raw[curr:hierPartEnd]
		authorityInfo, err = parseAuthority(hierPart)
		if err != nil {
			return nil, ErrInvalidURI
		}
		if hierPartEnd+1 < len(raw) {
			query = raw[hierPartEnd+1:]
		}
		curr = hierPartEnd + 1
	}

	if queryEnd == len(raw) {
		if hierPartEnd < 0 {
			hierPart = raw[curr:queryEnd]
			authorityInfo, err = parseAuthority(hierPart)
			if err != nil {
				return nil, ErrInvalidURI
			}
		} else {
			query = raw[curr:hierPartEnd]
			hierPart = raw[hierPartEnd+1 : queryEnd]
			authorityInfo, err = parseAuthority(hierPart)
			if err != nil {
				return nil, ErrInvalidURI
			}
		}
		u := &uri{
			scheme:    scheme,
			hierPart:  hierPart,
			authority: authorityInfo,
			query:     query,
		}
		return u, u.Validate()
	} else if queryEnd > 0 {
		if hierPartEnd < 0 {
			hierPart = raw[curr:queryEnd]
			authorityInfo, err = parseAuthority(hierPart)
			if err != nil {
				return nil, ErrInvalidURI
			}
		} else {
			query = raw[curr:queryEnd]
		}
		if queryEnd+1 < len(raw) {
			fragment = raw[queryEnd+1:]
		}
	}

	u := &uri{
		scheme:    scheme,
		hierPart:  hierPart,
		query:     query,
		fragment:  fragment,
		authority: authorityInfo,
	}

	return u, u.Validate()
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

var (
	rexScheme   = regexp.MustCompile(`^[\p{L}][\p{L}\d\+-\.]+$`)
	rexFragment = regexp.MustCompile(`^([\p{L}\d\-\._~\:@!\$\&'\(\)\*\+,;=\?/]|(%[[:xdigit:]]{2})+)+$`)
	rexQuery    = rexFragment
	rexSegment  = regexp.MustCompile(`^([\p{L}\d\-\._~\:@!\$\&'\(\)\*\+,;=]|(%[[:xdigit:]]{2})+)+$`)
	rexHostname = regexp.MustCompile(`^[\p{L}](([-\p{L}\d]+)?[\p{L}\d])?(\.\p{L}(([-\p{L}\d]+)?[\p{L}\d])?)*$`)

	// unreserved | pct-encoded | sub-delims
	rexRegname = regexp.MustCompile(`^([\p{L}\d\-\._~!\$\&'\(\)\*\+,;=]|(%[[:xdigit:]]{2})+)+$`)
	// unreserved | pct-encoded | sub-delims | ":"
	rexUserInfo = regexp.MustCompile(`^([\p{L}\d\-\._~\:!\$\&'\(\)\*\+,;=\?/]|(%[[:xdigit:]]{2})+)+$`)

	rexIPv6Zone = regexp.MustCompile(`:+[^%]?%\d+([\w\d]+)?$`)
	rexPort     = regexp.MustCompile(`^\d+$`)
)

// Validate checks that all parts of a URI abide by allowed characters
func (u *uri) Validate() error {
	if u.scheme != "" {
		if ok := rexScheme.MatchString(u.scheme); !ok {
			return ErrInvalidScheme
		}
	}
	if u.query != "" {
		if ok := rexQuery.MatchString(u.query); !ok {
			return ErrInvalidQuery
		}
	}
	if u.fragment != "" {
		if ok := rexFragment.MatchString(u.fragment); !ok {
			return ErrInvalidFragment
		}
	}
	if u.hierPart != "" {
		if u.authority != nil {
			a := u.Authority()
			if a != nil {
				return a.Validate(u.scheme)
			}
		}
	}
	return nil
}

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
func (a authorityInfo) Path() string     { return a.path }
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

func (a authorityInfo) Validate(schemes ...string) error {
	for _, segment := range strings.Split(a.path, "/") {
		if segment == "" {
			continue
		}
		if ok := rexSegment.MatchString(segment); !ok {
			return ErrInvalidPath
		}
	}

	if a.host != "" {
		if ok := rexIPv6Zone.MatchString(a.host); ok {
			z := strings.Index(a.host, "%")
			a.host = a.host[0:z]
		}
		isIP := net.ParseIP(a.host) != nil
		if !isIP {
			var isHost bool
			for _, scheme := range schemes {

				switch scheme {
				// for a few schemes with obvious reference to registered name
				// as a dns name, use hostname validation rather than RFC3986
				// registerer name (more accurate)
				// this allows for instance, to detect that 256.256.256.256 is an invalid ip address.
				//
				// see https://www.iana.org/assignments/uri-schemes/uri-schemes.xhtml
				// and feel free to add more
				case "http", "https", "dns", "telnet", "ldap", "dntp", "git", "ftp", "sftp", "smtp", "snmp", "imap",
					"svn", "udp", "tcp", "vnc", "ws", "soap", "rtsp", "rmi", "nntp", "ntp", "nfs", "mailto", "jms",
					"irc", "finger":
					isHost = rexHostname.MatchString(a.host)
				default:
					isHost = rexRegname.MatchString(a.host)
				}
				if !isHost {
					return ErrInvalidHost
				}
			}
		}
	}

	if a.port != "" {
		if ok := rexPort.MatchString(a.port); !ok {
			return ErrInvalidPort
		}
	}

	if a.userinfo != "" {
		if ok := rexUserInfo.MatchString(a.userinfo); !ok {
			return ErrInvalidUserInfo
		}
	}

	return nil
}

var (
	// byte literals
	atBytes       = []byte("@")
	colonBytes    = []byte(":")
	fragmentBytes = []byte("#")
	queryBytes    = []byte("?")
	slashBytes    = []byte("/")
)

func parseAuthority(hier string) (*authorityInfo, error) {
	// as per RFC 3986 Section 3.6
	var prefix, userinfo, host, port, path string

	// authority sections can begin with a '//'
	if strings.HasPrefix(hier, "//") {
		prefix = "//"
		hier = strings.TrimPrefix(hier, "//")
	}

	slashEnd := strings.Index(hier, "/")
	if prefix == "" {
		path = hier
	} else {
		// authority   = [ userinfo "@" ] host [ ":" port ]
		if slashEnd > 0 {
			if slashEnd < len(hier) {
				path = hier[slashEnd:]
			}
			hier = hier[:slashEnd]
		}

		host = hier
		if at := strings.Index(host, "@"); at > 0 {
			userinfo = host[:at]
			if at+1 < len(host) {
				host = host[at+1:]
			}
		}

		if bracket := strings.Index(host, "["); bracket >= 0 {
			// ipv6 addresses: "[" xx:yy:zz "]":port
			rawHost := host
			closingbracket := strings.Index(host, "]")
			if closingbracket > 0 {
				host = host[bracket+1 : closingbracket-bracket]
				rawHost = rawHost[closingbracket+1:]
			} else {
				return nil, ErrInvalidURI
			}
			if colon := strings.Index(rawHost, ":"); colon >= 0 {
				if colon+1 < len(rawHost) {
					port = rawHost[colon+1:]
				}
			}
		} else {
			if colon := strings.Index(host, ":"); colon > 0 {
				if colon+1 < len(host) {
					port = host[colon+1:]
				}
				host = host[:colon]
			}
		}
	}

	return &authorityInfo{
		prefix:   prefix,
		userinfo: userinfo,
		host:     host,
		port:     port,
		path:     path,
	}, nil
}

func (u *uri) ensureAuthorityExists() {
	if u.authority == nil {
		u.authority = &authorityInfo{}
	}
}

func (u *uri) SetScheme(scheme string) Builder {
	u.scheme = scheme
	return u
}

func (u *uri) SetUserInfo(userinfo string) Builder {
	u.ensureAuthorityExists()
	u.authority.userinfo = userinfo
	return u
}

func (u *uri) SetHost(host string) Builder {
	u.ensureAuthorityExists()
	u.authority.host = host
	return u
}

func (u *uri) SetPort(port string) Builder {
	u.ensureAuthorityExists()
	u.authority.port = port
	return u
}

func (u *uri) SetPath(path string) Builder {
	u.ensureAuthorityExists()
	u.authority.path = path
	return u
}

func (u *uri) SetQuery(query string) Builder {
	u.query = query
	return u
}

func (u *uri) SetFragment(fragment string) Builder {
	u.fragment = fragment
	return u
}

func (u *uri) Build() URI {
	return u
}

func (u *uri) Builder() Builder {
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
