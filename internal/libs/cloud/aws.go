package cloud

import "github.com/aws/aws-lambda-go/lambda"

func StartLambda(handler interface{}) {
	lambda.Start(handler)
}
