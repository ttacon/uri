package uri

import (
	"bytes"
	"fmt"
)

// MySQLURIBuilder builds a URI for use with the gomysql database/sql
// driver. Useage is fairly simple:
//
//  b := NewMySQLURIBuilder()
//  b.SetUser(user)
//  b.SetPassword(password)
//  b.SetHost(host)
//  b.SetPort(port)
//
//  driverURI := b.String()
type MySQLURIBuilder interface {
	URIBuilder

	// mysql specific fun
	SetUser(user string) MySQLURIBuilder
	SetPassword(password string) MySQLURIBuilder
	SetDB(db string) MySQLURIBuilder
}

type mysqlURI struct {
	*uri
	username string
	password string
	db       string
}

// NewMySQLURIBuilder returns a URI builder for builder
// gomysql driver URIs.
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

func (m *mysqlURI) SetDB(db string) MySQLURIBuilder {
	m.db = db
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
		"%s:%s@tcp(%s:%s)/%s?",
		m.username,
		m.password,
		host,
		port,
		m.db,
	)
}

// PostgresqlURIBuilder is a URI builder for
// connecting to Postgresql databases. Usage is similar to
// MySQLURIBuilder.
type PostgresqlURIBuilder interface {
	URIBuilder

	// mysql specific fun
	SetUser(user string) PostgresqlURIBuilder
	SetPassword(password string) PostgresqlURIBuilder
	SetDB(db string) PostgresqlURIBuilder
}

type postgresqlURI struct {
	*uri
	username string
	password string
	db       string
}

func NewPostgresqlURIBuilder() PostgresqlURIBuilder {
	return &postgresqlURI{
		uri: &uri{
			authority: &authorityInfo{},
		},
	}
}

func (m *postgresqlURI) SetUser(user string) PostgresqlURIBuilder {
	m.username = user
	return m
}

func (m *postgresqlURI) SetPassword(password string) PostgresqlURIBuilder {
	m.password = password
	return m
}

func (m *postgresqlURI) SetDB(db string) PostgresqlURIBuilder {
	m.db = db
	return m
}

func (m *postgresqlURI) String() string {
	// This is how postgresql connection strings should look:
	//
	// postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]

	var buf = bytes.NewBufferString("postgresql://")
	if len(m.username) > 0 {
		buf.WriteString(m.username)
		if len(m.password) > 0 {
			buf.Write(colonBytes)
			buf.WriteString(m.password)
		}

		buf.Write(atBytes)
	}

	if len(m.uri.authority.host) > 0 {
		buf.WriteString(m.uri.authority.host)
	}

	if len(m.uri.authority.port) > 0 {
		buf.Write(colonBytes)
		buf.WriteString(m.uri.authority.port)
	}

	if len(m.db) > 0 {
		buf.Write(slashBytes)
		buf.WriteString(m.db)
	}

	// TODO(ttacon): actually use options/query

	return buf.String()
}
