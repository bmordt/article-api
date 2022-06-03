# article-api
technical test article api

Assumptions
 - Docker is installed - for the spinning up a postgres instance
 - Golang is installed - I used go 1.16

To run tests:
`go test ./... -coverprofile=c.out`
To view coverprofile:
`go tool cover -html=c.out`

