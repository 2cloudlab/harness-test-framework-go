package main

import (
    "context"
	"fmt"
	lambda_context "github.com/aws/aws-lambda-go/lambda"
)

func LambdaHandler(ctx context.Context) (int, error) {
    fmt.Println("Hello World!")
    return 0, nil
}

func main() {
    lambda_context.Start(LambdaHandler)
}
