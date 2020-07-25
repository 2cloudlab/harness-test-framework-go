package main

import (
	"bytes"
	"context"
	"fmt"
	"os"

	lambda_context "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Performancer interface {
	Init()
	Start(ctx context.Context, params EventParams) []byte
}

var performer = S3Performancer{}

func Record(key string, value []byte) {
	input := &s3.PutObjectInput{
		Body:   bytes.NewReader(value),
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(key),
	}

	_, err := g_s3_service.PutObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
}

func LambdaHandler(ctx context.Context, params EventParams) (int, error) {
	lc, _ := lambdacontext.FromContext(ctx)
	Record(getReportName(params.RequestID, lc.AwsRequestID), performer.Start(ctx, params))
	return 0, nil
}

func main() {
	init_shared_resource()
	performer.Init()
	lambda_context.Start(LambdaHandler)
}
