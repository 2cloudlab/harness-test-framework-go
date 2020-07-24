package main

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func main() {
	svc := lambda.New(session.New())
	params := EventParams{Iteration: 5, LambdaFunctionName: "worker-handler", CountInSingleInstance: 1}
	payLoadInJson, _ := json.Marshal(params)
	input := &lambda.InvokeInput{
		FunctionName: aws.String("test-harness-framework"),
		Payload:      payLoadInJson,
	}
	_, err := svc.Invoke(input)
	if err != nil {
		recordError(err)
	}
}
