# FROM mcr.microsoft.com/devcontainers/typescript-node:22

# WORKDIR /workspace

# USER node


FROM mcr.microsoft.com/devcontainers/typescript-node:22 AS builder

WORKDIR /opt

ENV RYE_HOME="/opt/rye"
ENV PATH="$RYE_HOME/shims:$PATH"

# Install required packages
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        curl

# Install Rye and configure it
SHELL [ "/bin/bash", "-o", "pipefail", "-c" ]
RUN curl -sSf https://rye.astral.sh/get | RYE_INSTALL_OPTION="--yes" bash && \
    rye config --set-bool behavior.global-python=true && \
    rye config --set-bool behavior.use-uv=true

# Copy ALL project files, including LICENSE
COPY . ./

# Debug: List files to verify LICENSE exists
RUN ls -la

# Install dependencies with verbose output
RUN rye pin "$(cat .python-version)" && \
    rye sync -v

FROM mcr.microsoft.com/devcontainers/typescript-node:22
COPY --from=builder /opt/rye /opt/rye

WORKDIR /workspace

ENV RYE_HOME="/opt/rye"
ENV PATH="$RYE_HOME/shims:$PATH"
ENV PYTHONUNBUFFERED True

# Configure Rye for this stage
RUN rye config --set-bool behavior.global-python=true && \
    rye config --set-bool behavior.use-uv=true

# Set ownership to node user
RUN chown -R node $RYE_HOME

USER node