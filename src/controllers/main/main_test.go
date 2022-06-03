package main

import (
	"os"

	"testing"
)

func TestInitEnvVariables(t *testing.T) {
	testApiPort := "test-port"
	os.Setenv("APIPORT", testApiPort)
	testDBUser := "test-user"
	os.Setenv("DBUSER", testDBUser)
	testDBName := "test-db-name"
	os.Setenv("DBNAME", testDBName)
	testDBPassword := "test-db-password"
	os.Setenv("DBPASSWORD", testDBPassword)
	testDBHost := "test-db-host"
	os.Setenv("DBHOST", testDBHost)
	testDBPort := "test-db-port"
	os.Setenv("DBPORT", testDBPort)

	t.Run("Given all the env variables are set, no fatal log is thrown", func(t *testing.T) {
		initEnvVariables()
	})
}
