set CURDIR=%~dp0
set GOBIN=%CURDIR%\bin
go install -race ./...
pause
