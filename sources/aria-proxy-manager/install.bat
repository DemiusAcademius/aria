echo off
set GOARCH=amd64
set GOOS=linux
go build
rem mv aria-proxy-manager ..\..\content\aria-services\aria-proxy\aria-proxy-manager\docker-content\bin