# Installation

## Using Homebrew (MacOS/Linux)
If you have homebrew, the installation is as simple as:
```bash
brew tap kuzxnia/loadbot
brew install loadbot
```

## Running directly in Docker container

Alternatively, you can use the
[official Docker images](https://hub.docker.com/repository/docker/kuzxnia/loadbot)
that come with the `loadbot` binary pre-installed:
```bash
docker run --rm kuzxnia/loadbot --help
```

## By downloading the binaries (MacOS/Linux)

1. Go to the [releases](https://github.com/kuzxnia/loadbot/releases) and download
   the latest release archive for your platform.
2. Extract the archive.
3. Move the binary to somewhere in your `PATH`.

Sample steps for MacOS:
```bash
$ VERSION=<VERSION_TAG>
$ wget https://github.com/kuzxnia/loadbot/releases/download/${VERSION}/loadbot_${VERSION}_darwin_x86_64.tar.gz
$ tar -xvzf loadbot_${VERSION}_darwin_x86_64.tar.gz
$ mv loadbot /usr/local/bin
$ loadbot --help
```
