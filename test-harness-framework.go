package main

import (
	"context"
	"encoding/json"
	"fmt"

	lambda_context "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
)

// Each request will route to this function.
// So handle your request here.
// Also make sure your lambda function is idempotent.
// You should use fmt.Println() to log key info to Log Stream
func LambdaHandler(ctx context.Context, params EventParams) (string, error) {
	lc, _ := lambdacontext.FromContext(ctx)
	params.RequestID = lc.AwsRequestID
	payLoadInJson, _ := json.Marshal(params)
	// Log payload to Log Stream for better analyze this function behavior
	fmt.Println(payLoadInJson)
	input := &lambda.InvokeInput{
		FunctionName: aws.String(params.LambdaFunctionName),
		// Make the invoke async
		InvocationType: aws.String("Event"),
		Payload:        payLoadInJson,
	}
	for i := 0; i < params.NumberOfTasks; i++ {
		_, err := g_lambda_service.Invoke(input)
		if err != nil {
			recordError(err)
		}
	}
	// Return requestID and success flag
	return lc.AwsRequestID, nil
}

// Lambda will call main function only one time, when initialize.
// So you should initialize global resources, such as external services client, in the main entry.
func main() {
	init_shared_resource()
	// The following command will block and wait for the request sending from clients.
	// After idle for a duration, Lambda will automatically destroy the function instance.
	lambda_context.Start(LambdaHandler)
	// The following command will not log out to AWS Log Stream.
	fmt.Println("End of this function instance")
}
