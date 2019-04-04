set CGO_ENABLED=0
set GOROOT_BOOTSTRAP=C:/Go
::x86块
set GOARCH=386
set GOOS=linux
go build main.go
ren main iotserv
set GOARCH=amd64
set GOOS=windows
pause