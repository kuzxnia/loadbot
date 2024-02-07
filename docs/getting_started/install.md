# Installation

## Using Homebrew (MacOS/Linux)
If you have homebrew, the installation is as simple as:
```bash
brew tap kuzxnia/loadbot
brew install lbot
```

## By downloading the binaries (MacOS/Linux)

1. Go to the [releases](https://github.com/utkuozdemir/pv-migrate/releases) and download
   the latest release archive for your platform.
2. Extract the archive.
3. Move the binary to somewhere in your `PATH`.

Sample steps for MacOS:
```bash
$ VERSION=<VERSION_TAG>
$ wget https://github.com/utkuozdemir/pv-migrate/releases/download/${VERSION}/pv-migrate_${VERSION}_darwin_x86_64.tar.gz
$ tar -xvzf pv-migrate_${VERSION}_darwin_x86_64.tar.gz
$ mv pv-migrate /usr/local/bin
$ pv-migrate --help
```


## Running directly in Docker container

Alternatively, you can use the
[official Docker images](https://hub.docker.com/repository/docker/utkuozdemir/pv-migrate)
that come with the `pv-migrate` binary pre-installed:
```bash
docker run --rm -it kuzxnia/lbot:<IMAGE_TAG> .
docker run --rm -it kuzxnia/lbot-agent:<IMAGE_TAG> .
```

## 

Using Homebrew (MacOS/Linux)
If you have homebrew, the installation is as simple as:

brew tap utkuozdemir/pv-migrate
brew install pv-migrate

* install brew /tap
brew tap kuzxnia/loadbot
brew install lbot

* build from sources

* build with docker

