
build-docker: build-docker-lbot build-docker-agent

build-docker-lbot:
	docker build -t lbot -f docker/Dockerfile.lbot .
build-docker-agent:
	docker build -t lbot-agent -f docker/Dockerfile.agent .

run-docker-lbot:
	docker run --net=host -t --rm lbot
run-docker-agent:
	docker run -p 1234:1234 -t --rm lbot-agent

.PHONY: build-docker build-docker-agent build-docker-lbot run-docker-lbot run-docker-agent
