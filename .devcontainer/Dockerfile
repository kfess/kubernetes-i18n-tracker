FROM golang:1.25-bookworm

# Install system dependencies including Node.js and Python prerequisites
RUN apt-get update && apt-get install -y \
    git \
    make \
    curl \
    vim \
    bash-completion \
    ca-certificates \
    gnupg \
    python3 \
    python3-pip \
    python3-venv \
    && rm -rf /var/lib/apt/lists/*

# Install Node.js (LTS version)
RUN curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - \
    && apt-get install -y nodejs \
    && rm -rf /var/lib/apt/lists/*

# Install Yarn
RUN npm install -g yarn

ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=1000

RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME -s /bin/bash

USER $USERNAME

# ユーザーのホームディレクトリにGOPATHを設定
ENV GOPATH=/home/vscode/go
ENV PATH=$PATH:$GOPATH/bin:/home/vscode/.local/bin

# Install Python packages for data analysis
RUN pip3 install --user --no-cache-dir --break-system-packages \
    jupyter \
    pandas \
    matplotlib \
    plotly \
    numpy \
    seaborn

RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && \
    go install github.com/air-verse/air@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install golang.org/x/tools/gopls@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /workspace

ENV GOCACHE=/workspace/.cache/go-build
ENV GOMODCACHE=/workspace/.cache/go-mod