FROM golang:1.25-bookworm

RUN apt-get update && apt-get install -y \
    git \
    make \
    curl \
    vim \
    && rm -rf /var/lib/apt/lists/*

RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && \
    go install github.com/air-verse/air@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install golang.org/x/tools/gopls@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /workspace

ENV GOCACHE=/go/cache
ENV GOMODCACHE=/go/pkg/mod

ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
    && mkdir -p /go/cache /go/pkg/mod \
    && chown -R $USERNAME:$USERNAME /go

USER $USERNAME