module github.com/OpenIoTHub/server-go

go 1.12

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20191011191535-87dc89f01550
	golang.org/x/net => github.com/golang/net v0.0.0-20191011234655-491137f69257
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys => github.com/golang/sys v0.0.0-20191010194322-b09406accb47
	golang.org/x/text => github.com/golang/text v0.3.2
	golang.org/x/tools => github.com/golang/tools v0.0.0-20191011211836-4c025a95b26e
	golang.org/x/xerrors => github.com/golang/xerrors v0.0.0-20191011141410-1b5146add898
)

require (
	github.com/OpenIoTHub/utils v0.0.0-20200419073446-65f594ffa254
	github.com/iotdevice/zeroconf v0.0.0-20190527085138-7225942b5495 // indirect
	github.com/templexxx/cpufeat v0.0.0-20180724012125-cef66df7f161 // indirect
	github.com/templexxx/xor v0.0.0-20181023030647-4e92f724b73b // indirect
	github.com/urfave/cli/v2 v2.1.1
	github.com/xtaci/kcp-go v5.4.11+incompatible
	gopkg.in/yaml.v2 v2.2.4
)
