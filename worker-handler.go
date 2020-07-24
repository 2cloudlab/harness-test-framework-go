package main

import (
	"context"
	"fmt"
	"time"

	lambda_context "github.com/aws/aws-lambda-go/lambda"
)

func LambdaHandler(ctx context.Context, params EventParams) (int, error) {
	for i := 0; i < params.CountInSingleInstance; i++ {
		time.Sleep(2 * time.Second)
		fmt.Println("Hello World!")
	}
	return 0, nil
}

func main() {
	lambda_context.Start(LambdaHandler)
}
