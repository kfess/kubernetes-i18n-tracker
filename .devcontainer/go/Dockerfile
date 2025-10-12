FROM golang:1.25-bookworm

RUN apt-get update && apt-get install -y \
    git \
    make \
    curl \
    vim \
    bash-completion \
    && rm -rf /var/lib/apt/lists/*

RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && \
    go install github.com/air-verse/air@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install golang.org/x/tools/gopls@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /workspace

ENV GOCACHE=/workspace/.cache/go-build
ENV GOMODCACHE=/workspace/.cache/go-mod

ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=1000

RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME -s /bin/bash

USER $USERNAME