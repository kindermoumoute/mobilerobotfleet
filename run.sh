#!/usr/bin/env bash
`GOOS=linux go build`
docker-compose.exe up &
read -p "Press enter to continue"
docker-compose.exe down