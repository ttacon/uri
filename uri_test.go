package uri

import (
	"reflect"
	"testing"
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
				authorityInfo{
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
				authorityInfo{
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
				authorityInfo{
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
				authorityInfo{
					"user",
					"git.openstack.org",
					"29418",
					"/openstack/keystone.git",
				},
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
			t.Errorf("got back unexpected uri: %v != %v", got, test.uri)
		}
	}
}
