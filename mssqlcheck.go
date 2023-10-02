package main

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"net/url"
)

func IsMSSQLServer(uri *url.URL) string {
	sqldb, err := sql.Open("sqlserver",
		"sqlserver://"+uri.User.String()+"@"+uri.Host+"?database="+uri.Query().Get("db"))
	if err != nil {
		return "connect failed"
	}
	defer func() {
		_ = sqldb.Close()
	}()
	err = sqldb.Ping()
	if err != nil {
		return "ping failed"
	}
	var rows *sql.Rows
	rows, err = sqldb.Query(uri.Query().Get("q"))
	if err != nil {
		return "query failed"
	}
	cnt := 0
	for rows.Next() {
		var val any
		err = rows.Scan(&val)
		if err != nil {
			break
		}
		cnt++
	}
	if cnt == 0 {
		return "no rows returned"
	}

	return "verified"
}
