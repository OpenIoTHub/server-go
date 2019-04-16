rm -rf %GODIST%/natcloud/snap/server/*
set CGO_ENABLED=0
set GOROOT_BOOTSTRAP=C:/Go
::x86块
set GOOS=linux

set GOARCH=386
go build -ldflags -w main.go
ren main serverLinux386
upx serverLinux386
mv serverLinux386 %GODIST%/natcloud/snap/server/
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
  
::x64块
set GOARCH=amd64
go build -ldflags -w main.go
ren main serverLinuxAMD64
upx serverLinuxAMD64
mv serverLinuxAMD64 %GODIST%/natcloud/snap/server/
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
  
::arm块
set GOARCH=arm
go build -ldflags -w main.go
ren main serverLinuxArm
upx serverLinuxArm
mv serverLinuxArm %GODIST%/natcloud/snap/server/
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
  
::mips块
set GOARCH=mips64le
go build -ldflags -w main.go
ren main serverLinuxMips64le
upx serverLinuxMips64le
mv serverLinuxMips64le %GODIST%/natcloud/snap/server/

set GOARCH=mips64
go build -ldflags -w main.go
ren main serverLinuxMips64
upx serverLinuxMips64
mv serverLinuxMips64 %GODIST%/natcloud/snap/server/

set GOARCH=mipsle
set CGO_ENABLED=0
set GOMIPS=softfloat
go build -ldflags -w main.go
ren main serverLinuxMipsle
upx serverLinuxMipsle
mv serverLinuxMipsle %GODIST%/natcloud/snap/server/
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
set GOARCH=amd64
set GOOS=windows
pause