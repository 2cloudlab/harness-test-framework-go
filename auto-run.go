package main

import (
	"bytes"
	"encoding/json"
	"time"
	"sort"
	"strings"
	"io/ioutil"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
)

func upload() {
	svc := s3.New(session.New())
	
	// generate 16 objects with size range from 1KB to 32 MB, increase by a factor of 2 
	initSizeInBytes := 1024
	bucket_name := "2cloudlab-performance-benchmark-bucket"
	for i := 0; i < 16; i++ {
		subKey := getObjectName(i + 1)
		_, err := svc.HeadObject(&s3.HeadObjectInput{
			Bucket: aws.String(bucket_name),
			Key:    aws.String(subKey),
		})
		
		if err != nil {
			//could not find object
			input := &s3.PutObjectInput{
				Body:   bytes.NewReader(make([]byte, initSizeInBytes )),
				Bucket: aws.String(bucket_name),
				Key:    aws.String(subKey),
			}
			_, err := svc.PutObject(input)

			if err != nil {
				recordError(err)
			}
		}
		initSizeInBytes *= 2
	}
}

func generate_report(prefix []byte) {
	// Creating the array for JSON
	m := []interface{}{}
	jsonInStr := "[{\"abc\":100,\"xyz\":-2.1},{\"abc\":1000,\"xyz\":10.1}]"
    // Parsing/Unmarshalling JSON encoding/json
	err := json.Unmarshal([]byte(jsonInStr), &m)
	if err != nil {
        panic(err)
	}
	headersToMap := map[string][]float64{}
	headers := []string{}
	for _, obj := range m {
		for key, _ := range obj.(map[string]interface{}) {
			headersToMap[key] = []float64{}
			headers = append(headers, key)
		}
		break
	}

	for _, obj := range m {
		for key, val := range obj.(map[string]interface{}) {
			headersToMap[key] = append(headersToMap[key], val.(float64))
		}
	}

	for _, key := range headers {
		sort.Sort(sort.Reverse(sort.Float64Slice(headersToMap[key])))
	}
	var buffer strings.Builder
	buffer.WriteString(strings.Join(headers[:], " "))
	buffer.WriteString("\n")

	flat_data := []float64{}
	record_number := 0
	headers_number := len(headers)
	for _, key := range headers {
		flat_data = append(flat_data, headersToMap[key]...)
		record_number = len(headersToMap[key])
	}

	for i := 0; i < record_number; i++ {
		one_row := []float64{}
		for j := 0; j < headers_number; j++ {
			one_row = append(one_row, flat_data[i + j * record_number])
		}
		buffer.WriteString(strings.Trim(fmt.Sprint(one_row), "[]"))
		buffer.WriteString("\n")
	} 

	d1 := []byte(strings.Trim(buffer.String(),"\n"))
    ioutil.WriteFile("report.csv", d1, 0644)
}

func main() {
	// upload data to S3
	upload()
	// launch Lambda Function
	svc := lambda.New(session.New())
	params := EventParams{Iteration: 5, LambdaFunctionName: "worker-handler", CountInSingleInstance: 1}
	payLoadInJson, _ := json.Marshal(params)
	input := &lambda.InvokeInput{
		FunctionName: aws.String("test-harness-framework"),
		Payload:      payLoadInJson,
	}
	result, err := svc.Invoke(input)
	if err != nil {
		recordError(err)
		return
	}

	//wait 15 minutes
	time.Sleep(15 * time.Minute)

	//generate report
	generate_report(result.Payload)
}
