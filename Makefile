args = `arg="$(filter-out $@,$(MAKECMDGOALS))" && echo $${arg:-${1}}`

CONFIG_FILE=



build-docker: build-docker-lbot build-docker-agent

build-docker-lbot:
	docker build -t lbot -f docker/Dockerfile.lbot .
build-docker-agent:
	docker build -t lbot-agent -f docker/Dockerfile.agent .

build-dev:
	# only your platform
	goreleaser build --snapshot --clean --single-target

build-all:
	# only your platform
	goreleaser build --snapshot --clean

run-docker-lbot:
	@docker run --net=host -t --rm lbot $(call args,)
run-docker-lbot-config:
	docker run --net=host -i --rm lbot config --stdin < $(CONFIG_FILE)

run-docker-agent:
	docker run -p 1234:1234 -t --rm lbot-agent


NEXT_VERSION=$(shell svu next)
release:
	git tag -a $(NEXT_VERSION) -m "Release $(NEXT_VERSION)"


.PHONY: build-docker build-docker-agent build-docker-lbot run-docker-lbot run-docker-agent
