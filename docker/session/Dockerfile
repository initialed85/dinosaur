FROM ubuntu:22.04 AS builder

RUN apt-get update && \
    apt-get install -y \
    curl git entr \
    procps file strace screen \
    net-tools inetutils-ping traceroute netcat tcpdump iproute2

# java
RUN apt-get update && \
    apt-get install -y default-jdk

# lua
RUN apt-get update && \
    apt-get install -y lua5.4 luarocks

# python3
RUN apt-get update && \
    apt-get install -y python3 python3-dev python3-pip

# node
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.3/install.sh | bash && \
    . /root/.bashrc && nvm alias default v20 && nvm install v20 && nvm use v20

# typescript
RUN apt-get update && \
    apt-get install -y npm
RUN npm install -g typescript ts-node

# rust
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs > install.sh && \
  chmod +x install.sh && \
  ./install.sh -y
ENV PATH=${PATH}:/root/.cargo/bin/
RUN rustup update

# go
RUN curl -LO https://go.dev/dl/go1.21.3.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.21.3.linux-amd64.tar.gz
ENV PATH=${PATH}:/usr/local/go/bin

RUN CGO_ENABLED=0 go install github.com/sorenisanerd/gotty@latest

RUN mkdir -p /srv/cmd/

WORKDIR /srv/

# because Lua has no support for sockets
RUN luarocks install luasocket

# because TypeScript needs to know about some NodeJS concepts
RUN npm i --save-dev @types/node

RUN echo 'termcapinfo xterm* ti@:te@' >> /root/.screenrc

# this just turns it all into a single layer with the hope that it'll provide a slight image size optimisation
FROM --platform=linux/amd64 ubuntu:22.04

COPY --from=builder / /

WORKDIR /srv/

ENV PATH=${PATH}:/root/.cargo/bin/
ENV PATH=${PATH}:/usr/local/go/bin
ENV PATH=${PATH}:/root/go/bin

COPY docker-entrypoint.sh /docker-entrypoint.sh
COPY loop.sh /loop.sh
COPY build.sh /build.sh
COPY run.sh /run.sh
COPY watch.sh /watch.sh

ENV GOTTY_PORT=${GOTTY_PORT:-8080}
ENV GOTTY_PATH=${GOTTY_PATH}
ENV BASE_FOLDER_PATH=${BASE_FOLDER_PATH:-/srv/}
ENV BUILD_CMD=${BUILD_CMD}
ENV RUN_CMD=${RUN_CMD}

ENTRYPOINT ["/docker-entrypoint.sh"]
