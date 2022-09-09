FROM ubuntu:22.04

RUN apt-get update && \
    apt-get install -y \
    curl git npm entr \
    golang-1.18 python3 default-jdk lua5.4 luarocks \
    procps file strace \
    net-tools inetutils-ping traceroute netcat tcpdump iproute2

RUN npm install -g typescript ts-node

RUN curl https://sh.rustup.rs -sSf | sh -s -- -y

RUN luarocks install luasocket

ENV PATH=${PATH}:/usr/lib/go-1.18/bin:/root/go/bin/:/root/.cargo/bin

RUN go install github.com/sorenisanerd/gotty@latest

RUN mkdir -p /srv/cmd/

WORKDIR /srv/

COPY docker-entrypoint.sh /docker-entrypoint.sh
COPY loop.sh /loop.sh
COPY build.sh /build.sh
COPY run.sh /run.sh

ENV GOTTY_PORT=${GOTTY_PORT:-8080}
ENV GOTTY_PATH=${GOTTY_PATH}
ENV BASE_FOLDER_PATH=${BASE_FOLDER_PATH:-/srv/}
ENV BUILD_CMD=${BUILD_CMD}
ENV RUN_CMD=${RUN_CMD}

ENTRYPOINT ["/docker-entrypoint.sh"]
