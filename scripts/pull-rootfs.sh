#!/usr/bin/env bash

set -euo pipefail

ALPINE_VERSION="3.23.4"
ARCH="x86_64"

ARCHIVE="alpine-minirootfs-${ALPINE_VERSION}-${ARCH}.tar.gz"
URL="https://dl-cdn.alpinelinux.org/alpine/v3.23/releases/${ARCH}/${ARCHIVE}"

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ROOTFS="${PROJECT_ROOT}/rootfs"

echo "Creating rootfs directory..."
rm -rf "${ROOTFS}"
mkdir -p "${ROOTFS}"

echo "Downloading Alpine Mini RootFS..."
curl -L "${URL}" -o "/tmp/${ARCHIVE}"

echo "Extracting..."
tar -xzf "/tmp/${ARCHIVE}" -C "${ROOTFS}"

echo "Cleaning up..."
rm "/tmp/${ARCHIVE}"

echo "Done!"
echo "Root filesystem extracted to: ${ROOTFS}"