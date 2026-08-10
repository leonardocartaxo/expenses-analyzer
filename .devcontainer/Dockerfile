FROM node:24-bookworm-slim

# Slim glibc image: Alpine/musl breaks JetBrains WebStorm Remote Dev (host-status).
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
  && chmod 0440 /etc/sudoers.d/node
