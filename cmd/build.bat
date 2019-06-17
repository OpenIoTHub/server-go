set CGO_ENABLED=0
set GOROOT_BOOTSTRAP=C:/Go
::x86块
set GOARCH=386
set GOOS=windows
go build -ldflags -w main.go
ren main.exe windows386.exe
::upx windows386.exe
set GOOS=linux
go build -ldflags -w main.go
ren main linux386
upx linux386
set GOOS=freebsd
go build -ldflags -w main.go
ren main freebsd386
upx freebsd386
set GOOS=darwin
go build -ldflags -w main.go
ren main darwin386
upx darwin386
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
  
::x64块
set GOARCH=amd64
set GOOS=windows
go build -ldflags -w main.go
ren main.exe windowsAmd64.exe
::upx windowsAmd64.exe
set GOOS=linux
go build -ldflags -w main.go
ren main linuxAMD64
upx linuxAMD64
set GOOS=freebsd
go build -ldflags -w main.go
ren main freebsdAMD64
upx freebsdAMD64
set GOOS=darwin
go build -ldflags -w main.go
ren main darwinAMD64
upx darwinAMD64
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
  
::arm块
set GOARCH=arm
set GOOS=linux
go build -ldflags -w main.go
ren main LinuxArm
upx LinuxArm
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
  
::mips块
set GOARCH=mips64le
set GOOS=linux
go build -ldflags -w main.go
ren main LinuxMips64le
upx LinuxMips64le

set GOARCH=mips64
set GOOS=linux
go build -ldflags -w main.go
ren main LinuxMips64
upx LinuxMips64

set GOARCH=mipsle
set GOOS=linux
set CGO_ENABLED=0
set GOMIPS=softfloat
go build -ldflags -w main.go
ren main LinuxMipsle
upx LinuxMipsle

set GOARCH=mips
set GOOS=linux
set CGO_ENABLED=0
set GOMIPS=softfloat
go build -ldflags -w main.go
ren main LinuxMips
upx LinuxMips
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
set GOARCH=amd64
set GOOS=windows
pause