FROM node:24-bookworm-slim

# Slim glibc image: Alpine/musl breaks JetBrains WebStorm Remote Dev (host-status).
# Docker CLI is client-only (Docker-outside-of-Docker via mounted docker.sock).
ARG DOCKER_VERSION=27.5.1
ARG COMPOSE_VERSION=v2.32.4

RUN apt-get update \
  && apt-get install -y --no-install-recommends \
    bash \
    ca-certificates \
    curl \
    fontconfig \
    git \
    gh \
    libfreetype6 \
    libxi6 \
    libxext6 \
    libxrender1 \
    libxtst6 \
    openssh-client \
    procps \
    sudo \
    zsh \
  && rm -rf /var/lib/apt/lists/* \
  && echo "node ALL=(root) NOPASSWD:ALL" > /etc/sudoers.d/node \
  && chmod 0440 /etc/sudoers.d/node \
  && arch="$(uname -m)" \
  && case "${arch}" in \
       x86_64) darch=x86_64; carch=x86_64 ;; \
       aarch64|arm64) darch=aarch64; carch=aarch64 ;; \
       *) darch="${arch}"; carch="${arch}" ;; \
     esac \
  && curl -fsSL "https://download.docker.com/linux/static/stable/${darch}/docker-${DOCKER_VERSION}.tgz" \
    | tar -xz -C /tmp \
  && install -m 0755 /tmp/docker/docker /usr/local/bin/docker \
  && rm -rf /tmp/docker \
  && mkdir -p /usr/local/lib/docker/cli-plugins \
  && curl -fsSL "https://github.com/docker/compose/releases/download/${COMPOSE_VERSION}/docker-compose-linux-${carch}" \
    -o /usr/local/lib/docker/cli-plugins/docker-compose \
  && chmod +x /usr/local/lib/docker/cli-plugins/docker-compose
