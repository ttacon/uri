package uri

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_rawURIParse(t *testing.T) {
	var tests = []struct {
		uriRaw string
		uri    *uri
		err    error
	}{
		{
			"foo://example.com:8042/over/there?name=ferret#nose",
			&uri{"foo", "//example.com:8042/over/there", "name=ferret", "nose",
				&authorityInfo{
					"//",
					"",
					"example.com",
					"8042",
					"/over/there",
				},
			},
			nil,
		},
		{
			"http://httpbin.org/get?utf8=%e2%98%83",
			&uri{"http", "//httpbin.org/get", "utf8=%e2%98%83", "",
				&authorityInfo{
					"//",
					"",
					"httpbin.org",
					"",
					"/get",
				},
			},
			nil,
		},
		{
			"mailto://user@domain.com",
			&uri{"mailto", "//user@domain.com", "", "",
				&authorityInfo{
					"//",
					"user",
					"domain.com",
					"",
					"",
				},
			},
			nil,
		},
		{
			"ssh://user@git.openstack.org:29418/openstack/keystone.git",
			&uri{"ssh", "//user@git.openstack.org:29418/openstack/keystone.git", "", "",
				&authorityInfo{
					"//",
					"user",
					"git.openstack.org",
					"29418",
					"/openstack/keystone.git",
				},
			},
			nil,
		},
		{
			"https://willo.io/#yolo",
			&uri{"https", "//willo.io/", "", "yolo",
				&authorityInfo{"//", "", "willo.io", "", "/"},
			},
			nil,
		},
	}

	for _, test := range tests {
		got, err := Parse(test.uriRaw)
		if err != test.err {
			t.Errorf("got back unexpected err: %v != %v", err, test.err)
			continue
		} else if !reflect.DeepEqual(got, test.uri) {
			t.Errorf("got back unexpected (raw: %s), uri: %v != %v",
				test.uriRaw, fmt.Sprintf("%#v", got), fmt.Sprintf("%#v", test.uri))
		}
	}
}

func Test_rawURIParseFailed(t *testing.T) {
	var tests = []struct {
		uriRaw string
		uri    *uri
		err    error
	}{
		{
			"http://httpbin.org/get?utf8=\xe2\x98\x83",
			&uri{"http", "//httpbin.org/get", "utf8=\xe2\x98\x83", "",
				&authorityInfo{
					"//",
					"",
					"httpbin.org",
					"",
					"/get",
				},
			},
			ErrInvalidQuery,
		},
		{
			// without // prefix, this is a path!
			"mailto:user@domain.com",
			&uri{"mailto", "user@domain.com", "", "",
				&authorityInfo{
					path: "user@domain.com",
				},
			},
			nil,
		},
	}

	for _, test := range tests {
		got, err := Parse(test.uriRaw)
		if err != test.err {
			t.Errorf("got back unexpected err: %v != %v", err, test.err)
			continue
		} else if !reflect.DeepEqual(got, test.uri) {
			t.Errorf("got back unexpected (raw: %s), uri: %v != %v",
				test.uriRaw, fmt.Sprintf("%#v", got), fmt.Sprintf("%#v", test.uri))
		}
	}
}

func Test_ParseThenString(t *testing.T) {
	var tests = []string{
		"foo://example.com:8042/over/there?name=ferret#nose",
		"http://httpbin.org/get?utf8=yödeléï",
		"http://httpbin.org/get?utf8=%e2%98%83",
		"mailto://user@domain.com",
		"ssh://user@git.openstack.org:29418/openstack/keystone.git",
		"https://willo.io/#yolo",
	}

	for _, test := range tests {
		uri, err := Parse(test)
		if err != nil {
			t.Errorf("failed to parse URI %q, err: %v", test, err)
		} else if uri.String() != test {
			t.Errorf("uri.String() != test: %v != %v", uri.String(), test)
		}
	}
	_, err := Parse("http://httpbin.org/get?utf8=\xe2\x98\x83")
	// no normalization takes place at the moment
	assert.Error(t, err)
}

func Benchmark_Parse(b *testing.B) {
	var tests = []string{
		"foo://example.com:8042/over/there?name=ferret#nose",
		"http://httpbin.org/get?utf8=%e2%98%83",
		"mailto://user@domain.com",
		"ssh://user@git.openstack.org:29418/openstack/keystone.git",
		"https://willo.io/#yolo",
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Parse(tests[i%5])
	}
}

func Benchmark_String(b *testing.B) {
	var tests = []*uri{
		{"foo", "//example.com:8042/over/there", "name=ferret", "nose",
			&authorityInfo{"//", "", "example.com", "8042", "/over/there"},
		},
		{"http", "//httpbin.org/get", "utf8=\xe2\x98\x83", "",
			&authorityInfo{"//", "", "httpbin.org", "", "/get"},
		},
		{"mailto", "user@domain.com", "", "",
			&authorityInfo{"//", "user", "domain.com", "", ""},
		},
		{"ssh", "//user@git.openstack.org:29418/openstack/keystone.git", "", "",
			&authorityInfo{"//", "user", "git.openstack.org", "29418", "/openstack/keystone.git"},
		},
		{"https", "//willo.io/", "", "yolo",
			&authorityInfo{"//", "", "willo.io", "", "/"},
		},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = tests[i%5].String()
	}
}

func Test_Building(t *testing.T) {
	var tests = []struct {
		uri, uriChanged string
		name            string
	}{
		{
			"mailto://user@domain.com",
			"http://yolo@newdomain.com:443",
			"yolo",
		},
		{
			"https://user@domain.com",
			"http://yolo2@newdomain.com:443",
			"yolo2",
		},
	}

	for _, test := range tests {
		auri, err := Parse(test.uri)
		if err != nil {
			t.Errorf("failed to parse uri: %v", err)
			continue
		}
		nuri := auri.Builder().SetUserInfo(test.name).SetHost("newdomain.com").SetScheme("http").SetPort("443")
		zuri, ok := nuri.(URI)
		assert.True(t, ok)
		assert.Equal(t, "//"+test.name+"@newdomain.com:443", zuri.Authority().String())
		assert.Equal(t, "443", nuri.URI().Authority().Port())
		val := nuri.String()
		if val != test.uriChanged {
			t.Logf("val: %#v", val)
			t.Logf("test: %#v", test.uriChanged)
			t.Errorf("vals don't match: %v != %v", val, test.uriChanged)
		}
		assert.Equal(t, "http", nuri.URI().Scheme())

		_ = nuri.SetPath("/abcd")
		assert.Equal(t, "/abcd", nuri.URI().Authority().Path())

		_ = nuri.SetQuery("a=b&x=5").SetFragment("chapter")
		assert.Equal(t, url.Values{"a": []string{"b"}, "x": []string{"5"}}, nuri.URI().Query())
		assert.Equal(t, "chapter", nuri.URI().Fragment())
		assert.Equal(t, test.uriChanged+"/abcd?a=b&x=5#chapter", nuri.URI().String())
		assert.Equal(t, test.uriChanged+"/abcd?a=b&x=5#chapter", nuri.String())
	}

	// build from scratch
	u, _ := Parse("http:")
	b := u.Builder()
	//uri_test.go:251: ""

	nu := (b.URI()).(*uri)
	assert.Nil(t, nu.authority)
	assert.Equal(t, "", u.Authority().UserInfo())

	b = b.SetUserInfo("user:pwd").SetHost("newdomain").SetPort("444")
	assert.Equal(t, "http://user:pwd@newdomain:444", b.String())
}

// TestMoreURI borrows from other URI validators to exercise strict RFC3986
// conformance (taken from .Net, perl, python, )
func TestMoreURI(t *testing.T) {
	invalidURIs := []string{
		"mailto://{}:{}@host.domain.com",
		"https://user:passwd@[FF02::3::5:8080",
		"https://user:passwd@[FF02::3::5:8080/?#",
		"https://user:passwd@[FF02::3::5:8080#",
		"https://user:passwd@[FF02::3::5:8080#abc",

		// this test comes from the format test in JSONSchema-test suite
		"//foo.bar/?baz=qux#quux", // missing scheme and //

		// from https://docs.microsoft.com/en-gb/dotnet/api/system.uri.iswellformeduristring?view=netframework-4.7.2#System_Uri_IsWellFormedUriString_System_String_System_UriKind_
		"http://www.contoso.com/path???/file name", // The string is not correctly escaped.
		"c:\\directory\filename",                   // The string is an absolute Uri that represents an implicit file Uri.
		"http:\\host/path/file",                    // The string contains unescaped backslashes even if they will be treated as forward slashes
		"www.contoso.com/path/file",                // The string represents a hierarchical absolute Uri and does not contain "://"
		"2013.05.29_14:33:41",                      // relative URIs with a colon (':') in their first segment are not considered well-formed.

		// from https://metacpan.org/source/SONNEN/Data-Validate-URI-0.07/t/is_uri.t
		"",
		"foo",
		"foo@bar",
		"http://<foo>",      // illegal characters
		"://bob/",           // empty scheme
		"1http://bob",       // bad scheme
		"http:////foo.html", // bad path
		"http://example.w3.org/%illegal.html",
		"http://example.w3.org/%a",     // partial escape
		"http://example.w3.org/%a/foo", // partial escape
		"http://example.w3.org/%at",    // partial escape

		// from https://github.com/python-hyper/rfc3986/blob/master/tests/test_validators.py
		"https://user:passwd@[21DA:00D3:0000:2F3B:02AA:00FF:FE28:9C5A]:8080:8090/a?query=value#fragment", // multiple ports
		"https://user:passwd@[FF02::3::5]:8080/a?query=value#fragment",                                   // invalid IPv6
		"https://user:passwd@[FADF:01%en0]:8080/a?query=value#fragment",                                  // invalid IPv6
		"https://user:passwd@256.256.256.256:8080/a?query=value#fragment",                                // invalid IPv4
		"https://user:passwd@[FADF:01%en0:8080/a?query=value#fragment",                                   // invalid IPv6 (missing bracket)

		// from github.com/scalatra/rl: URI parser in scala
		"http://www.exa mple.org",

		// and others..
		"?invalidscheme://www.example.com",
		"inv;alidscheme://www.example.com",
		"http://www.example.org/hello/world.txt/?id=5&pa{}rt=three#there-you-go", // invalid char in query
		"http://www.example.org/hello/world.txt/?id=5&part=three#there-you-go{}", // invalid char in fragment
	}
	validURIs := []string{
		"urn://example-bin.org/path",
		"https://example-bin.org/path",
		"https://example-bin.org/path?",
		"mailto://u:p@host.domain.com#",  // empty fragment
		"mailto://u:p@host.domain.com?#", // empty query + fragment
		"http:",
		"foo:",

		// this one is dubious: Microsoft (.Net) recognize the C:/... string as a path and
		// states this as incorrect uri -- all other validators state a host "c" and state this uri as a valid one
		"file://c:/directory/filename",

		// from https://metacpan.org/source/SONNEN/Data-Validate-URI-0.07/t/is_uri.t
		// (many of those come from the rfc3986 examples)
		"http://localhost/",
		"http://example.w3.org/path%20with%20spaces.html",
		"http://example.w3.org/%20",
		"ftp://ftp.is.co.za/rfc/rfc1808.txt",
		"ftp://ftp.is.co.za/../../../rfc/rfc1808.txt",
		"http://www.ietf.org/rfc/rfc2396.txt",
		"ldap://[2001:db8::7]/c=GB?objectClass?one",
		"mailto:John.Doe@example.com",   // valid but counter-intuitive: userinfo is actually a path
		"mailto://John.Doe@example.com", // this is the right way
		"news:comp.infosystems.www.servers.unix",
		"tel:+1-816-555-1212",
		"telnet://192.0.2.16:80/",
		"urn:oasis:names:specification:docbook:dtd:xml:4.1.2",
		"http://www.richardsonnen.com/",

		// from https://github.com/python-hyper/rfc3986/blob/master/tests/test_validators.py
		"ssh://ssh@git.openstack.org:22/sigmavirus24",
		"https://git.openstack.org:443/sigmavirus24",
		"ssh://git.openstack.org:22/sigmavirus24?foo=bar#fragment",
		"git://github.com",
		"https://user:passwd@[21DA:00D3:0000:2F3B:02AA:00FF:FE28:9C5A]:8080/a?query=value#fragment",
		"https://user:passwd@[::1%25lo]:8080/a?query=value#fragment",
		"https://user:passwd@[FF02:30:0:0:0:0:0:5%25en1]:8080/a?query=value#fragment",
		"https://user:passwd@127.0.0.1:8080/a?query=value#fragment",
		"https://user:passwd@http-bin.org:8080/a?query=value#fragment",

		// from github.com/scalatra/rl: URI parser in scala
		"http://www.example.org:8080",
		"http://www.example.org/",
		"http://www.詹姆斯.org/",
		"http://www.example.org/hello/world.txt",
		"http://www.example.org/hello/world.txt/?id=5&part=three",
		"http://www.example.org/hello/world.txt/?id=5&part=three#there-you-go",
		"http://www.example.org/hello/world.txt/#here-we-are",
	}

	for _, invURI := range invalidURIs {
		res := IsURI(invURI)
		if assert.Falsef(t, res, "expected %q to be an invalid URI", invURI) {
			t.Logf("%q is an invalid URI as expected", invURI)
		}
	}
	for _, validURI := range validURIs {
		res := IsURI(validURI)
		assert.True(t, res, "expected %q to be a valid URI", validURI)
	}
}

func Test_MoreParse(t *testing.T) {
	_, err := Parse("1http://bob")
	assert.Equal(t, ErrInvalidScheme, err)

	_, err = Parse("http://www.example.org/hello/world.txt/?id=5&part=three#there-you-go")
	assert.NoError(t, err)

	_, err = Parse("http://www.example.org/hélloô/mötor/world.txt/?id=5&part=three#there-you-go")
	assert.NoError(t, err)

	_, err = Parse("http://www.example.org/hello/world.txt/?id=5&pa{}rt=three#there-you-go")
	assert.Equal(t, ErrInvalidQuery, err)

	_, err = Parse("http://www.example.org/hello/yzx;=1.1/world.txt/?id=5&part=three#there-you-go")
	assert.NoError(t, err)

	_, err = Parse("http://www.example.org/hello/{}yzx;=1.1/world.txt/?id=5&part=three#there-you-go")
	assert.Equal(t, ErrInvalidPath, err)

	_, err = Parse("https://user:passwd@127.0.0.1:8080/a?query=value#fragment")
	assert.NoError(t, err)

	_, err = Parse("https://user:passwd@286;0.0.1:8080/a?query=value#fragment")
	assert.Equal(t, ErrInvalidHost, err)

	_, err = Parse("http://www.詹姆斯.org/")
	assert.NoError(t, err)

	_, err = Parse("https://user:passwd@256.256.256.256:8080/a?query=value#fragment")
	assert.Equal(t, ErrInvalidHost, err)

	_, err = Parse("file://c:/directory/filename")
	assert.NoError(t, err)

	_, err = Parse("ldap://[2001:db8::7]/c=GB?objectClass?one")
	assert.NoError(t, err)

	_, err = Parse("ldap://[2001:db8::7]:8080/c=GB?objectClass?one")
	assert.NoError(t, err)

	_, err = Parse("https://user:passwd@[FF02:30:0:0:0:0:0:5%25]:8080/a?query=value#fragment")
	assert.NoError(t, err)

	_, err = Parse("https://user:passwd@[FF02:30:0:0:0:0:0:5%25en0]:8080/a?query=value#fragment")
	assert.NoError(t, err)

	_, err = Parse("https://user:passwd@[FF02:30:0:0:0:0:0:5%25en0:8080/a?query=value#fragment") // lack closing bracket
	assert.Equal(t, ErrInvalidURI, err)

	_, err = Parse("https://user:passwd@[FF02:30:0:0:0:0:0:5%25lo]:8080/a?query=value#fragment")
	assert.NoError(t, err)

	_, err = Parse("https://user:passwd@[21DA:00D3:0000:2F3B:02AA:00FF:FE28:9C5A]:8080:8090/a?query=value#fragment")
	assert.Equal(t, ErrInvalidPort, err)

	_, err = Parse("tel:+1-816-555-1212")
	assert.NoError(t, err)

	_, err = Parse("http+unix://%2Fvar%2Frun%2Fsocket/path?key=value")
	assert.NoError(t, err)

	_, err = Parse("https://user{}:passwd@[FF02:30:0:0:0:0:0:5%25en0]:8080/a?query=value#fragment")
	assert.Equal(t, ErrInvalidUserInfo, err)

	u, err := Parse("urn:oasis:names:specification:docbook:dtd:xml:4.1.2")
	assert.NoError(t, err)
	assert.Equal(t, "urn", u.Scheme())

	_, err = Parse("https://user:passwd@[21DA:00D3:0000:2F3B:02AA:00FF:FE28:9C5A]:8080/a?query=value#fragment")
	assert.NoError(t, err)

	u, err = Parse("https://example-bin.org/path?")
	assert.NoError(t, err)
	assert.Equal(t, "/path", u.Authority().Path())

	u, err = Parse("https://example-bin.org/path#frag?withQuestionMark")
	assert.NoError(t, err)
	assert.Equal(t, "/path", u.Authority().Path())
	nuri := u.(*uri)
	assert.Equal(t, "", nuri.query)
	assert.Equal(t, "frag?withQuestionMark", u.Fragment())

	u, err = Parse("mailto://u:p@host.domain.com?#")
	assert.NoError(t, err)
	assert.Equal(t, "", u.Authority().Path())

	u, err = Parse("mailto://u:p@host.domain.com?#ahahah")
	assert.NoError(t, err)
	assert.Equal(t, "", u.Authority().Path())
	nuri = u.(*uri)
	assert.Equal(t, "", nuri.query)
	assert.Equal(t, "ahahah", u.Fragment())

	u, err = Parse("ldap://[2001:db8::7]/c=GB?objectClass?one")
	assert.NoError(t, err)
	assert.Equal(t, "/c=GB", u.Authority().Path())
	nuri = u.(*uri)
	assert.Equal(t, "objectClass?one", nuri.query)
	assert.Equal(t, "", u.Fragment())

	u, err = Parse("http://www.example.org/hello/world.txt/?id=5&part=three")
	assert.NoError(t, err)
	assert.Equal(t, "/hello/world.txt/", u.Authority().Path())
	nuri = u.(*uri)
	assert.Equal(t, "id=5&part=three", nuri.query)
	assert.Equal(t, "", u.Fragment())

	u, err = Parse("http://www.example.org/hello/world.txt/?id=5&part=three?another#abc?efg")
	assert.NoError(t, err)
	assert.Equal(t, "/hello/world.txt/", u.Authority().Path())
	nuri = u.(*uri)
	assert.Equal(t, "id=5&part=three?another", nuri.query)
	assert.Equal(t, "abc?efg", u.Fragment())
	assert.Equal(t, url.Values{"id": []string{"5"}, "part": []string{"three?another"}}, u.Query())

	u, err = Parse("?")
	assert.Error(t, err)

	u, err = Parse("#")
	assert.Error(t, err)

	u, err = Parse("?#")
	assert.Error(t, err)

	u, err = Parse("")
	assert.Error(t, err)

	u, err = Parse(" ")
	assert.Error(t, err)

	u, err = Parse("ht?tp://host")
	assert.Error(t, err)

	u, err = Parse("https://user:passwd@[21DA:00D3:0000:2F3B:02AA:00FF:FE28:9C5A%25en0]:8080/a?query=value#fragment")
	assert.NoError(t, err)
	assert.Equal(t, "21DA:00D3:0000:2F3B:02AA:00FF:FE28:9C5A%25en0", u.Authority().Host())
	assert.Equal(t, "//user:passwd@[21DA:00D3:0000:2F3B:02AA:00FF:FE28:9C5A%25en0]:8080/a", u.Authority().String())
	assert.Equal(t, "https", u.Scheme())
	assert.Equal(t, url.Values{"query": []string{"value"}}, u.Query())
}

func Test_Edge(t *testing.T) {
	u, err := Parse("https:")
	assert.NoError(t, err)
	assert.Equal(t, "https", u.Scheme())

	u, err = Parse("ht?tps:")
	assert.Error(t, err)

	u, err = Parse("https://user:passwd@[21DA:00D3:0000:2F3B:02AA:00FF:FE28:9C5A%25]:8080/a?query=value#fragment")
	assert.NoError(t, err)
	assert.Equal(t, "21DA:00D3:0000:2F3B:02AA:00FF:FE28:9C5A%25", u.Authority().Host())

	u, err = Parse("https://user:passwd@[::1%25lo]:8080/a?query=value#fragment")
	assert.NoError(t, err)
	assert.Equal(t, "https", u.Scheme())
	assert.Equal(t, "8080", u.Authority().Port())
	assert.Equal(t, "user:passwd", u.Authority().UserInfo())

	// empty host
	_, err = Parse("https://user:passwd@:8080/a?query=value#fragment")
	assert.Equal(t, ErrMissingHost, err)
}
