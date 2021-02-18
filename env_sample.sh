#!/bin/bash -e

# github
export GITHUB_USER="your-name"
export GITHUB_TOKEN="your-token"

# GO
export GOPRIVATE="github.com/CyberAgent/*"

# DB
export DB_MASTER_HOST="db"
export DB_MASTER_USER="hoge"
export DB_MASTER_PASSWORD="moge"
export DB_SLAVE_HOST="db"
export DB_SLAVE_USER="hoge"
export DB_SLAVE_PASSWORD="moge"
export DB_LOG_MODE="true"

# aws
export AWS_REGION="ap-northeast-1"

# grpc server
export IAM_SVC_ADDR="iam:8002"
export FINDING_SVC_ADDR="finding:8001"

# notification alert url
export NOTIFICATION_ALERT_URL="http://localhost:8080/#/alert/alert"
