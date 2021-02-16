#!/bin/bash -e

cd "$(dirname "$0")"

# load env
. ../env.sh

# setting remote repository
TAG="local-test-$(date '+%Y%m%d')"
IMAGE_FINDING="core/finding"
IMAGE_IAM="core/iam"
IMAGE_PROJECT="core/project"
IMAGE_ALERT="core/alert"
IMAGE_REPORT="core/report"
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query "Account" --output text)
REGISTORY="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com"

# build & push
aws ecr get-login-password --region ${AWS_REGION} \
  | docker login \
    --username AWS \
    --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com

docker build --build-arg GITHUB_USER=${GITHUB_USER} --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} -t ${IMAGE_FINDING}:${TAG} ../src/finding/
docker build --build-arg GITHUB_USER=${GITHUB_USER} --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} -t ${IMAGE_IAM}:${TAG} ../src/iam/
docker build --build-arg GITHUB_USER=${GITHUB_USER} --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} -t ${IMAGE_PROJECT}:${TAG} ../src/project/
docker build --build-arg GITHUB_USER=${GITHUB_USER} --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} -t ${IMAGE_ALERT}:${TAG} ../src/alert/
docker build --build-arg GITHUB_USER=${GITHUB_USER} --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} -t ${IMAGE_REPORT}:${TAG} ../src/report/

docker tag ${IMAGE_FINDING}:${TAG} ${REGISTORY}/${IMAGE_FINDING}:${TAG}
docker tag ${IMAGE_IAM}:${TAG}     ${REGISTORY}/${IMAGE_IAM}:${TAG}
docker tag ${IMAGE_PROJECT}:${TAG} ${REGISTORY}/${IMAGE_PROJECT}:${TAG}
docker tag ${IMAGE_ALERT}:${TAG}   ${REGISTORY}/${IMAGE_ALERT}:${TAG}
docker tag ${IMAGE_REPORT}:${TAG}  ${REGISTORY}/${IMAGE_REPORT}:${TAG}

docker push ${REGISTORY}/${IMAGE_FINDING}:${TAG}
docker push ${REGISTORY}/${IMAGE_IAM}:${TAG}
docker push ${REGISTORY}/${IMAGE_PROJECT}:${TAG}
docker push ${REGISTORY}/${IMAGE_ALERT}:${TAG}
docker push ${REGISTORY}/${IMAGE_REPORT}:${TAG}
