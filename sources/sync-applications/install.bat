@echo off
set GOARCH=amd64
set GOOS=linux
go build
move sync-applications ../../content/scripts/sync-applications