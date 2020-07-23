package main

import (
	"context"
	"encoding/json"
	"fmt"
	lambda_context "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func recordError(err error) {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case lambda.ErrCodeServiceException:
			fmt.Println(lambda.ErrCodeServiceException, aerr.Error())
		case lambda.ErrCodeResourceNotFoundException:
			fmt.Println(lambda.ErrCodeResourceNotFoundException, aerr.Error())
		case lambda.ErrCodeInvalidRequestContentException:
			fmt.Println(lambda.ErrCodeInvalidRequestContentException, aerr.Error())
		case lambda.ErrCodeRequestTooLargeException:
			fmt.Println(lambda.ErrCodeRequestTooLargeException, aerr.Error())
		case lambda.ErrCodeUnsupportedMediaTypeException:
			fmt.Println(lambda.ErrCodeUnsupportedMediaTypeException, aerr.Error())
		case lambda.ErrCodeTooManyRequestsException:
			fmt.Println(lambda.ErrCodeTooManyRequestsException, aerr.Error())
		case lambda.ErrCodeInvalidParameterValueException:
			fmt.Println(lambda.ErrCodeInvalidParameterValueException, aerr.Error())
		case lambda.ErrCodeEC2UnexpectedException:
			fmt.Println(lambda.ErrCodeEC2UnexpectedException, aerr.Error())
		case lambda.ErrCodeSubnetIPAddressLimitReachedException:
			fmt.Println(lambda.ErrCodeSubnetIPAddressLimitReachedException, aerr.Error())
		case lambda.ErrCodeENILimitReachedException:
			fmt.Println(lambda.ErrCodeENILimitReachedException, aerr.Error())
		case lambda.ErrCodeEFSMountConnectivityException:
			fmt.Println(lambda.ErrCodeEFSMountConnectivityException, aerr.Error())
		case lambda.ErrCodeEFSMountFailureException:
			fmt.Println(lambda.ErrCodeEFSMountFailureException, aerr.Error())
		case lambda.ErrCodeEFSMountTimeoutException:
			fmt.Println(lambda.ErrCodeEFSMountTimeoutException, aerr.Error())
		case lambda.ErrCodeEFSIOException:
			fmt.Println(lambda.ErrCodeEFSIOException, aerr.Error())
		case lambda.ErrCodeEC2ThrottledException:
			fmt.Println(lambda.ErrCodeEC2ThrottledException, aerr.Error())
		case lambda.ErrCodeEC2AccessDeniedException:
			fmt.Println(lambda.ErrCodeEC2AccessDeniedException, aerr.Error())
		case lambda.ErrCodeInvalidSubnetIDException:
			fmt.Println(lambda.ErrCodeInvalidSubnetIDException, aerr.Error())
		case lambda.ErrCodeInvalidSecurityGroupIDException:
			fmt.Println(lambda.ErrCodeInvalidSecurityGroupIDException, aerr.Error())
		case lambda.ErrCodeInvalidZipFileException:
			fmt.Println(lambda.ErrCodeInvalidZipFileException, aerr.Error())
		case lambda.ErrCodeKMSDisabledException:
			fmt.Println(lambda.ErrCodeKMSDisabledException, aerr.Error())
		case lambda.ErrCodeKMSInvalidStateException:
			fmt.Println(lambda.ErrCodeKMSInvalidStateException, aerr.Error())
		case lambda.ErrCodeKMSAccessDeniedException:
			fmt.Println(lambda.ErrCodeKMSAccessDeniedException, aerr.Error())
		case lambda.ErrCodeKMSNotFoundException:
			fmt.Println(lambda.ErrCodeKMSNotFoundException, aerr.Error())
		case lambda.ErrCodeInvalidRuntimeException:
			fmt.Println(lambda.ErrCodeInvalidRuntimeException, aerr.Error())
		case lambda.ErrCodeResourceConflictException:
			fmt.Println(lambda.ErrCodeResourceConflictException, aerr.Error())
		case lambda.ErrCodeResourceNotReadyException:
			fmt.Println(lambda.ErrCodeResourceNotReadyException, aerr.Error())
		default:
			fmt.Println(aerr.Error())
		}
	} else {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
	}
}

func LambdaHandler(ctx context.Context, params EventParams) (int, error) {
	payLoadInJson, _ := json.Marshal(params)
	svc := lambda.New(session.New())
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
	return 0, nil
}

func main() {
	fmt.Println("Before Start")
	lambda_context.Start(LambdaHandler)
	fmt.Println("After Start")
}
