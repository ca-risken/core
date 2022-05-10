#! /bin/sh
if [ "${IMAGE_TAG}" = "" ]; then
  IMAGE_TAG=latest
fi
if [ "${IMAGE_PREFIX}" = "" ]; then
  IMAGE_PREFIX=default_prefix
fi
docker build ${BUILD_OPT} -t ${IMAGE_NAME}:${IMAGE_TAG} .
