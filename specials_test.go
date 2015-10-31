package uri

import "testing"

func Test_MySQLURIBuilding(t *testing.T) {
	var tests = []struct {
		user, passwd, host, port string
		expected                 string
	}{
		{"", "", "", "", ":@tcp(localhost:3306)/?"},
		{"ttacon", "yolo", "", "", "ttacon:yolo@tcp(localhost:3306)/?"},
		{"ttacon", "yolo", "", "3307", "ttacon:yolo@tcp(localhost:3307)/?"},
		{"ttacon", "yolo", "db.secret.aweso.me", "", "ttacon:yolo@tcp(db.secret.aweso.me:3306)/?"},
	}

	for _, test := range tests {
		b := NewMySQLURIBuilder()
		b.SetUser(test.user)
		b.SetPassword(test.passwd)
		b.SetHost(test.host)
		b.SetPort(test.port)

		if b.String() != test.expected {
			t.Errorf("values don't match: %v != %v", b.String(), test.expected)
		}
	}
}

func Test_PostgresqlURIBuilding(t *testing.T) {
	var tests = []struct {
		user, passwd, host, port string
		expected                 string
	}{
		{"", "", "", "", "postgresql://"},
		{"ttacon", "yolo", "", "", "postgresql://ttacon:yolo@"},
		{"ttacon", "yolo", "", "1234", "postgresql://ttacon:yolo@"},
		{"ttacon", "yolo", "db.secret.aweso.me", "", "postgresql://ttacon:yolo@db.secret.aweso.me"},
	}

	for _, test := range tests {
		b := NewPostgresqlURIBuilder()
		b.SetUser(test.user)
		b.SetPassword(test.passwd)
		b.SetHost(test.host)

		if b.String() != test.expected {
			t.Errorf("values don't match: %v != %v", b.String(), test.expected)
		}
	}
}
