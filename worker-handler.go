package main

import (
	"context"

	lambda_context "github.com/aws/aws-lambda-go/lambda"
)

type Performancer interface {
	Init()
	Start(ctx context.Context, params EventParams) []byte
}

func LambdaHandler(ctx context.Context, params EventParams) (int, error) {
	performer.Start(ctx, params)
	return 0, nil
}

var performer = S3Performancer{}

func main() {
	performer.Init()
	lambda_context.Start(LambdaHandler)
}
