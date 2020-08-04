package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type EventParams struct {
	RequestID              string `json:"RequestID"`
	LambdaFunctionName     string `json:"FunctionName"`
	TaskName               string `json:"TaskName"`
	NumberOfTasks          int    `json:"NumberOfTasks"`
	ConcurrencyForEachTask int    `json:"ConcurrencyForEachTask"`
	RawJson                string `json:"RawJson"`
	NumberOfSamples        int    `json:"NumberOfSamples"`
}

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

func getObjectName(i uint8) string {
	return fmt.Sprintf("test-data/object-%d-1024-KB", i)
}

func getObjectSize(i uint8) string {
	s := 1 << (i - 1)
	if i > 10 {
		return fmt.Sprintf("%d MB", s/1024)
	} else {
		return fmt.Sprintf("%d KB", s)
	}
}

func getObjectSizeInBytes(i uint8) int {
	s := 1 << (i - 1)
	return s * 1024
}

func getReportName(prefix string, key string) string {
	return fmt.Sprintf("%s/%s", prefix, key)
}

func downloadFile(bucket string, key string) []byte {
	fmt.Println("worker handler request id is ", strings.Split(key, "/")[1])
	buf := aws.NewWriteAtBuffer([]byte{})
	_, err := g_s3_downloader.Download(buf,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		recordError(err)
	}
	return buf.Bytes()
}

func downloadByPrefix(bucket string, prefix string) [][]byte {
	resp, err := g_s3_service.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	results := [][]byte{}
	if err != nil {
		recordError(err)
		return results
	}
	for _, item := range resp.Contents {
		results = append(results, downloadFile(bucket, *item.Key))
	}
	return results
}

func readFile(fileName string) ([]byte, error) {
	// Open an file by fieName
	jsonFile, err := os.Open(fileName)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}
	// Read our opened file as a byte array.
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}
	return byteValue, nil
}

var g_s3_service *s3.S3
var g_s3_downloader *s3manager.Downloader
var g_lambda_service *lambda.Lambda

func init_shared_resource() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	g_s3_service = s3.New(sess)
	g_s3_downloader = s3manager.NewDownloader(sess)
	g_lambda_service = lambda.New(sess)
}
