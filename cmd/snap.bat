rm -rf %GODIST%/natcloud/snap/bin/server/*
set CGO_ENABLED=0
set GOROOT_BOOTSTRAP=C:/Go
::x86块
set GOOS=linux

set GOARCH=386
go build -ldflags -w main.go
ren main serverLinux386
upx serverLinux386
mv serverLinux386 %GODIST%/natcloud/snap/bin/server/
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
  
::arm块
set GOARCH=arm
go build -ldflags -w main.go
ren main serverLinuxArm
upx serverLinuxArm
mv serverLinuxArm %GODIST%/natcloud/snap/bin/server/
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
set GOARCH=amd64
set GOOS=windows
pause