#!/usr/bin/env bash
# ./01_create_lambda.bash
set -e

function create_lambda() {
 # ! Absolute path to code
 local dir=$PWD/src/hello-world-v0/

 # * LocalStack localhost endpoint
 local endpoint=http://localhost:4566

 # * Lambda configuration
 local function_handler=index.handler
 local function_name=hello-world-v0
 local function_role=arn:aws:iam::000000000000:role/localstack-does-not-care
 local function_runtime=nodejs18.x

 aws --endpoint-url $endpoint lambda delete-function --function-name $function_name || true

 # * Create lambda in LocalStack
 aws --endpoint-url $endpoint lambda create-function \
  --code S3Bucket="hot-reload",S3Key="$dir" \
  --function-name $function_name \
  --handler $function_handler \
  --role $function_role \
  --runtime $function_runtime
}

create_lambda