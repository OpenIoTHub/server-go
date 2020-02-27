# Server for OpenIoTHub
[![Build Status](https://travis-ci.com/OpenIoTHub/server-go.svg?branch=master)](https://travis-ci.com/OpenIoTHub/server-go)
## OpenIoTHub 物联网服务器

### 如果你有自建转发服务端的需求则可以自建此服务端，服务器需要同时开启指定的TCP，UDP端口
#### 建好此服务器之后，[网关](https://github.com/OpenIoTHub/server-go/releases )就可以连接自建的服务器了！

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
