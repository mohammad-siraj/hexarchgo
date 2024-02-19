#!/usr/bin/env bash
# ./02_invoke_lambda.bash
set -e

function invoke_lambda() {
 # * LocalStack localhost endpoint
 local endpoint=http://localhost:4566

 # * Lambda configuration
 local function_name=test_lambda

 # * Invoke configuration
 local invoke_payload='{ "name": "siraj mohammad" ,"email":"siraj2siraj"}'
 local invoke_to_file=response.json
 local invoke_type=RequestResponse

 # * Invoke lambda from LocalStack
 awslocal lambda invoke \
  --cli-binary-format raw-in-base64-out \
  --function-name $function_name \
  --invocation-type $invoke_type \
  --payload "$invoke_payload" \
  $invoke_to_file
}

invoke_lambda