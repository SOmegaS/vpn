FROM ubuntu:latest AS base

SHELL ["/bin/sh", "-eu", "-c"]

WORKDIR /build

RUN apt update

RUN apt-get update && apt-get install -y \
    sudo \
    wget \
    tar \
    build-essential \
    iproute2 \
    && apt-get clean

RUN useradd -ms /bin/bash dockeruser

RUN echo 'dockeruser ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

USER dockeruser

WORKDIR /home/dockeruser

ENV GO_VERSION 1.23.1

RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm go${GO_VERSION}.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

RUN go version

COPY . ./

ARG PORT=12345

EXPOSE ${PORT}

#RUN [ -c /dev/net/tun ] || (echo "TUN device not available" && exit 1)

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-gnu-gcc go build -o vpn ~/cmd/client/main.go

CMD ["sudo", "./vpn"]

