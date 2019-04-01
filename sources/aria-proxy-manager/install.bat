@echo off
set GOARCH=amd64
set GOOS=linux
go build
move aria-proxy-manager ../../content/aria-services/aria-proxy/aria-proxy-manager/docker-content/bin