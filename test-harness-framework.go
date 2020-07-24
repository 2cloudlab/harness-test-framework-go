package main

import (
	"context"
	"encoding/json"
	"fmt"

	lambda_context "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func LambdaHandler(ctx context.Context, params EventParams) (string, error) {
	lc, _ := lambdacontext.FromContext(ctx)
	params.RequestID = lc.AwsRequestID
	svc := lambda.New(session.New())
	payLoadInJson, _ := json.Marshal(params)
	input := &lambda.InvokeInput{
		FunctionName:   aws.String(params.LambdaFunctionName),
		InvocationType: aws.String("Event"),
		Payload:        payLoadInJson,
	}
	for i := 0; i < params.Iteration; i++ {
		_, err := svc.Invoke(input)
		if err != nil {
			recordError(err)
		}
	}
	return lc.AwsRequestID, nil
}

func main() {
	fmt.Println("Before Start")
	lambda_context.Start(LambdaHandler)
	fmt.Println("After Start")
}
