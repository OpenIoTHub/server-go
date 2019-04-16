rm -rf %GODIST%/natcloud/server/*
set CGO_ENABLED=0
set GOROOT_BOOTSTRAP=C:/Go
::x86块
set GOARCH=386
set GOOS=windows
go build -ldflags -w main.go
ren main.exe serverWindows386.exe
::upx windows386.exe
mv serverWindows386.exe %GODIST%/natcloud/server/
set GOOS=linux
go build -ldflags -w main.go
ren main serverLinux386
upx serverLinux386
mv serverLinux386 %GODIST%/natcloud/server/
set GOOS=freebsd
go build -ldflags -w main.go
ren main serverFreebsd386
upx serverFreebsd386
mv serverFreebsd386 %GODIST%/natcloud/server/
set GOOS=darwin
go build -ldflags -w main.go
ren main serverDarwin386
upx serverDarwin386
mv serverDarwin386 %GODIST%/natcloud/server/
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
  
::x64块
set GOARCH=amd64
set GOOS=windows
go build -ldflags -w main.go
ren main.exe serverWindowsAmd64.exe
::upx windowsAmd64.exe
mv serverWindowsAmd64.exe %GODIST%/natcloud/server/

set GOOS=linux
go build -ldflags -w main.go
ren main serverLinuxAMD64
upx serverLinuxAMD64
mv serverLinuxAMD64 %GODIST%/natcloud/server/

set GOOS=freebsd
go build -ldflags -w main.go
ren main serverFreebsdAMD64
upx freebsdAMD64
mv freebsdAMD64 %GODIST%/natcloud/server/

set GOOS=darwin
go build -ldflags -w main.go
ren main serverDarwinAMD64
upx serverDarwinAMD64
mv serverDarwinAMD64 %GODIST%/natcloud/server/
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
  
::arm块
set GOARCH=arm
set GOOS=linux
go build -ldflags -w main.go
ren main serverLinuxArm
upx serverLinuxArm
mv serverLinuxArm %GODIST%/natcloud/server/
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
  
::mips块
set GOARCH=mips64le
set GOOS=linux
go build -ldflags -w main.go
ren main serverLinuxMips64le
upx serverLinuxMips64le
mv serverLinuxMips64le %GODIST%/natcloud/server/

set GOARCH=mips64
set GOOS=linux
go build -ldflags -w main.go
ren main serverLinuxMips64
upx serverLinuxMips64
mv serverLinuxMips64 %GODIST%/natcloud/server/

set GOARCH=mipsle
set GOOS=linux
set CGO_ENABLED=0
set GOMIPS=softfloat
go build -ldflags -w main.go
ren main serverLinuxMipsle
upx serverLinuxMipsle
mv serverLinuxMipsle %GODIST%/natcloud/server/
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
set GOARCH=amd64
set GOOS=windows
pause