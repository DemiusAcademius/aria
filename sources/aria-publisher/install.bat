@echo off
set GOARCH=amd64
set GOOS=linux
go build
move aria-publisher ../../content/aria-services/aria-publisher/docker-content/bin