#!/bin/sh
if [[ "${ARCH}" = "" ]]; then
  echo "environment value ARCH is not set"
  exit 1
fi


if [[ "${ARCH}" = "arm64" ]]; then
  curl -OL https://go.dev/dl/go1.17.1.linux-arm64.tar.gz
  rm -rf /usr/local/go && tar -C /usr/local -xzf go1.17.1.linux-arm64.tar.gz
else
  goenv install 1.17.1
  goenv global 1.17.1
fi