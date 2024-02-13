
build-dev:
	# only your platform
	goreleaser build --snapshot --clean --single-target

build-all:
	# only your platform
	goreleaser build --snapshot --clean

.PHONY: build-dev build-all
