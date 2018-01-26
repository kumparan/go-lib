package sqlimporter

import "strings"

var mysqlDialect = map[string]string{
	"create": "CREATE DATABASE %s",
	"use":    "USE %s",
	"drop":   "DROP DATABASE %s",
}

var postgresDialect = map[string]string{
	"create": "CREATE SCHEMA %s",
	"use":    "SET search_path TO %s",
	"drop":   "DROP SCHEMA %s CASCADE",
}

func getDialect(driver, process string) string {
	switch strings.ToLower(driver) {
	case "mysql":
		return mysqlDialect[process]
	case "postgres":
		return postgresDialect[process]
	default:
		return mysqlDialect[process]
	}
}
