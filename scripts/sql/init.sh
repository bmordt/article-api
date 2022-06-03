#!/bin/bash

docker kill article-sql
docker rm article-sql

docker build -f Dockerfile -t article-sql:dev .

docker run -d -p 5432:5432 --name=article-sql -e POSTGRES_PASSWORD=12345 -e POSTGRES_DB=nine article-sql:dev