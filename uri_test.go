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
			&uri{
				"foo",
				"//example.com:8042/over/there",
				"name=ferret",
				"nose",
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
