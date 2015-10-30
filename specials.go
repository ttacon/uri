package uri

import "fmt"

type MySQLURIBuilder interface {
	URIBuilder

	// mysql specific fun
	SetUser(user string) MySQLURIBuilder
	SetPassword(password string) MySQLURIBuilder
}

type mysqlURI struct {
	*uri
	username string
	password string
}

func NewMySQLURIBuilder() MySQLURIBuilder {
	return &mysqlURI{
		uri: &uri{
			authority: &authorityInfo{},
		},
	}
}

func (m *mysqlURI) SetUser(user string) MySQLURIBuilder {
	m.username = user
	return m
}

func (m *mysqlURI) SetPassword(password string) MySQLURIBuilder {
	m.password = password
	return m
}

func (m *mysqlURI) String() string {
	host := m.uri.authority.host
	if len(host) == 0 {
		host = "localhost"
	}

	port := m.uri.authority.port
	if len(port) == 0 {
		port = "3306"
	}

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/",
		m.username,
		m.password,
		host,
		port,
	)
}
