package uri

import (
	"reflect"
	"testing"

	"github.com/ttacon/pretty"
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
			nil,
		},
		{
			"mailto:user@domain.com",
			&uri{"mailto", "user@domain.com", "", "",
				&authorityInfo{
					"",
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
		got, err := ParseURI(test.uriRaw)
		if err != test.err {
			t.Errorf("got back unexpected err: %v != %v", err, test.err)
			continue
		} else if !reflect.DeepEqual(got, test.uri) {
			t.Errorf("got back unexpected (raw: %s), uri: %v != %v",
				test.uriRaw, pretty.Sprintf("%#v", got), pretty.Sprintf("%#v", test.uri))
		}
	}
}

func Test_ParseThenString(t *testing.T) {
	var tests = []string{
		"foo://example.com:8042/over/there?name=ferret#nose",
		"http://httpbin.org/get?utf8=\xe2\x98\x83",
		"mailto:user@domain.com",
		"ssh://user@git.openstack.org:29418/openstack/keystone.git",
		"https://willo.io/#yolo",
	}

	for _, test := range tests {
		uri, err := ParseURI(test)
		if err != nil {
			t.Errorf("failed to parse URI, err: %v", err)
		} else if uri.String() != test {
			pretty.Println(uri)
			t.Errorf("uri.String() != test: %v != %v", uri.String(), test)
		}
	}
}

func Benchmark_Parse(b *testing.B) {
	var tests = []string{
		"foo://example.com:8042/over/there?name=ferret#nose",
		"http://httpbin.org/get?utf8=\xe2\x98\x83",
		"mailto:user@domain.com",
		"ssh://user@git.openstack.org:29418/openstack/keystone.git",
		"https://willo.io/#yolo",
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ParseURI(tests[i%5])
	}
}

func Benchmark_String(b *testing.B) {
	var tests = []*uri{
		&uri{"foo", "//example.com:8042/over/there", "name=ferret", "nose",
			&authorityInfo{"//", "", "example.com", "8042", "/over/there"},
		},
		&uri{"http", "//httpbin.org/get", "utf8=\xe2\x98\x83", "",
			&authorityInfo{"//", "", "httpbin.org", "", "/get"},
		},
		&uri{"mailto", "user@domain.com", "", "",
			&authorityInfo{"", "user", "domain.com", "", ""},
		},
		&uri{"ssh", "//user@git.openstack.org:29418/openstack/keystone.git", "", "",
			&authorityInfo{"//", "user", "git.openstack.org", "29418", "/openstack/keystone.git"},
		},
		&uri{"https", "//willo.io/", "", "yolo",
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
			"mailto:user@domain.com",
			"mailto:yolo@domain.com",
			"yolo",
		},
		{
			"https://user@domain.com",
			"https://yolo2@domain.com",
			"yolo2",
		},
	}

	for _, test := range tests {
		uri, err := ParseURI(test.uri)
		if err != nil {
			t.Errorf("failed to parse uri: %v", err)
			continue
		}

		val := uri.Builder().SetUserInfo(test.name).String()
		if val != test.uriChanged {
			t.Errorf("vals don't match: %v != %v", val, test.uriChanged)
		}
	}
}
