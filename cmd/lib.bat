set CGO_ENABLED=1
set GOROOT_BOOTSTRAP=C:/Go
::x86块
set GOARCH=386
set GOOS=linux
go build -buildmode=c-shared -o libiotserv.so main.go
set GOARCH=amd64
set GOOS=windows
pause