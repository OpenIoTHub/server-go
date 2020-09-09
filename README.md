# Server for OpenIoTHub
[![Build Status](https://travis-ci.com/OpenIoTHub/server-go.svg?branch=master)](https://travis-ci.com/OpenIoTHub/server-go)

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-white.svg)](https://snapcraft.io/server-go)
## OpenIoTHub 物联网服务器

### 如果你有自建转发服务端的需求则可以自建此服务端，服务器需要同时开启指定的TCP，UDP端口
#### 建好此服务器之后，[网关](https://github.com/OpenIoTHub/gateway-go/releases )就可以连接自建的服务器了！

```
如果使用redis持久化保存Http代理配置请将配置文件中的redis的使能打开enabled: true
```
```
如果使用redis后OpenIoTHub App无法添加http代理成功请尝试在redis-cli中执行config set stop-writes-on-bgsave-error no
参考：https://www.baidu.com/s?ie=UTF-8&wd=MISCONF Redis is configured to save RDB snapshots
127.0.0.1:6379> config set stop-writes-on-bgsave-error no
如果配置redis后正常使用请忽略本条问题
```

```yaml
my_public_ip_or_domian: "" #你运行本软件的服务器的域名或者ip地址，用来使用命令生成token时有用
common:
  bind_addr: 0.0.0.0 #服务器监听的IP地址，默认监听所有
  tcp_port: 34320 #服务器使用的tcp端口
  kcp_port: 34320 #服务器使用的kcp(实际上是udp)端口
  udp_p2p_port: 34321 #服务器使用的UDP端口，用于辅助p2p
  kcp_p2p_port: 34322 #服务器使用的kcp（UDP）端口，用于辅助p2p
  tls_port: 34321 #服务器使用的tls(实际上是tcp)端口，用于安全通信
  grpc_port: 34322 #服务器使用的grpc(实际上是tcp)端口，用于grpc通信
  http_port: 80 #服务器监听的http（tcp）端口，用于提供http代理功能
  https_port: 443 #服务器监听的https（tcp）端口，用于提供https代理功能
security:
  login_key: HLLdsa544&*S #用户自定义的服务器秘钥，此为默认，用户个人使用服务器请修改
  tls_Cert_file_path: ""
  tls_key_file_path: ""
  https_cert_file_path: ""
  https_key_file_path: ""
redisconfig:
  enabled: false #是否使用redis保存用户http代理配置 <----这里打开redis
  network: tcp  #redis使用tcp连接，默认即可
  address: 127.0.0.1:6379 #redis的地址，默认本机，redis默认端口6379，请根据自己的redis配置
  database: 0 #redis的默认服务器0，如果你不懂请保持0
  needAuth: false #redis是否需要密码验证，默认不需要false，如果你的redis需要密码请将false改为true并配置下面password为redis密码
  password: "" #redis的密码，needAuth:true时有效
```

You can install the pre-compiled binary (in several different ways),
use Docker.

Here are the steps for each of them:

## Install the pre-compiled binary

**homebrew tap** (only on macOS for now):

```sh
$ brew install OpenIoTHub/tap/server-go
```

**homebrew** (may not be the latest version):

```sh
$ brew install server-go（not support yet）
```

**snapcraft**:

```sh
$ sudo snap install server-go
```
*** 默认配置文件的路径：/root/snap/server-go/current/server-yaml

**scoop**:

```sh
$ scoop bucket add OpenIoTHub https://github.com/OpenIoTHub/scoop-bucket.git
$ scoop install server-go
```

**deb/rpm**:

Download the `.deb` or `.rpm` from the [releases page][releases] and
install with `dpkg -i` and `rpm -i` respectively.

**Shell script**:

```sh
$ curl -sfL https://install.goreleaser.com/github.com/OpenIoTHub/server-go.sh | sh
```

**manually**:

Download the pre-compiled binaries from the [releases page][releases] and
copy to the desired location.

## Running with Docker

You can also use it within a Docker container. To do that, you'll need to
execute something more-or-less like the following:

```sh
$ docker run openiothub/server:latest
```

Note that the image will almost always have the last stable Go version.

[releases]: https://github.com/OpenIoTHub/server-go/releases
