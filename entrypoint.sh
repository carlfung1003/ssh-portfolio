#!/bin/sh
# Generate host key on first boot (persisted on Fly volume)
KEY_DIR="${SSH_HOST_KEY%/*}"
mkdir -p "$KEY_DIR"
if [ ! -f "$SSH_HOST_KEY" ]; then
  echo "Generating SSH host key..."
  ssh-keygen -t ed25519 -f "$SSH_HOST_KEY" -N ""
fi
exec ssh-portfolio --serve
