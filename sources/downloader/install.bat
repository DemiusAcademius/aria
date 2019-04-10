@echo off
set GOARCH=amd64
set GOOS=windows
go build
move downloader.exe C:\Go\bin