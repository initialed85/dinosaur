FROM golang:1.21 AS build

WORKDIR /srv/

COPY go.mod /srv/go.mod
COPY go.sum /srv/go.sum
RUN --mount=type=cache,target=/root/.cache/go-build go mod download

COPY cmd /srv/cmd
COPY internal /srv/internal
COPY pkg /srv/pkg

RUN --mount=type=cache,target=/root/.cache/go-build go build -v -o main /srv/cmd/main.go

FROM ubuntu:22.04

RUN apt-get update && apt-get install -y ca-certificates curl gnupg lsb-release golang-1.18

RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
RUN echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

RUN apt-get update && apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

WORKDIR /srv/

RUN mkdir -p /srv/pkg/sessions/languages
COPY pkg/sessions/languages /srv/pkg/sessions/languages

RUN mkdir -p /srv/docker/
COPY docker/session /srv/docker/session

COPY --from=build /srv/main /srv/main

ENTRYPOINT ["bash", "-c", "docker network create dinosaur-external || true; docker network create dinosaur-internal --internal; /srv/main"]
CMD []
