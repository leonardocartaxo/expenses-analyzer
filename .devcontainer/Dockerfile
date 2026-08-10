FROM node:24-alpine

# Lean Alpine toolchain for Spec Kit + monorepo work (no app stack yet).
RUN apk add --no-cache \
    bash \
    ca-certificates \
    curl \
    git \
    github-cli \
    sudo \
    zsh \
  && echo "node ALL=(root) NOPASSWD:ALL" > /etc/sudoers.d/node \
  && chmod 0440 /etc/sudoers.d/node
