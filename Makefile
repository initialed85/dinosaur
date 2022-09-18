SHELL=/bin/bash

all: run

clean: down

build:
	docker build -t dinosaur-session -f docker/session/Dockerfile ./docker/session
	docker compose -f docker/docker-compose.yml build --parallel

up: build
	docker compose -f docker/docker-compose.yml up -d

logs:
	docker compose -f docker/docker-compose.yml logs -f -t

down:
	docker compose -f docker/docker-compose.yml down --remove-orphans --volumes || true

run: up
	# TODO: bit of a hack to support a teardown on Ctrl + C
	bash -c "trap '$(MAKE) down' EXIT; $(MAKE) logs"
