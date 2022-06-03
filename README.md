# article-api
technical test article api

Assumptions
 - Docker is installed - for the spinning up a postgres instance
 - Golang is installed - I used go 1.16

To run tests:
`go test ./... -coverprofile=c.out`
To view coverprofile:
`go tool cover -html=c.out`

To install all dependencies:
`go mod download`

The required Environment variables to run this API can be found in `.env.template`. These need to be set prior to running the API. This can be done simply as setting in an .env file and running: `export $(grep -v '^#' .env | xargs)`.
The DB variables are also set in the init.sh function when creating the new postgres instance
```
APIPORT - Port number the API will run on. e.g. 8080
DBUSER - Name of the postgres db user. e.g. postgres
DBNAME - article-sql
DBPASSWORD - password for the user to access the DB. e.g. 12345
DBHOST - Hostname of the DB. e.g. localhost
DBPORT - Port the DB is run on. e.g. 5432 for postgres
```

To run the api from the root directory: `go run src/controllers/main/main.go`
Also can be done from building the binary from the root dir: `go build ./src/controllers/main/main.go` and then `./main`

