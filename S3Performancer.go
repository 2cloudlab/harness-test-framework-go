package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3 performancer
type S3Performancer struct {
}

func (s3P S3Performancer) Start(ctx context.Context, params EventParams) map[string][]float64 {
	m := map[string]interface{}{}
	// The value of RawJson is specified in config.json.
	// If you don't need it, just do nothing on RawJson
	// We use RawJson to pass FileSize parameter, so the following code will parse RawJson and retrieve FileSize field.
	err := json.Unmarshal([]byte(params.RawJson), &m)
	if err != nil {
		recordError(err)
		return map[string][]float64{}
	}
	object_level := uint8(m["FileSize"].(float64))
	sample_data_key := getObjectName(object_level)

	// The value of ConcurrencyForEachTask is specified in config.json.
	// We use it to determine the nmumber of groutines
	operations := make(chan int, params.ConcurrencyForEachTask)
	// The value of NumberOfSamples is specified in config.json.
	// We use it to determine the number of operations that will be issued.
	operationsNumber := params.NumberOfSamples
	operationResults := make(chan time.Duration, operationsNumber)
	for g := 0; g < params.ConcurrencyForEachTask; g++ {
		go func(o <-chan int, results chan<- time.Duration) {
			for range o {
				benchmarkTimer := time.Now()
				result, err := g_s3_service.GetObject(&s3.GetObjectInput{
					Bucket: aws.String(os.Getenv("BUCKET_NAME")),
					Key:    aws.String(sample_data_key),
				})

				// if a request fails, exit
				if err != nil {
					recordError(err)
					panic("Failed to get object: " + err.Error())
				}

				buf := new(bytes.Buffer)
				buf.ReadFrom(result.Body)
				// stop the timer for this benchmark
				totalTime := time.Now().Sub(benchmarkTimer)

				results <- totalTime
			}
		}(operations, operationResults)
	}

	for i := 0; i < operationsNumber; i++ {
		operations <- i
	}

	close(operations)

	latencyInSeconds := []float64{}
	for s := 0; s < operationsNumber; s++ {
		l := <-operationResults
		latencyInSeconds = append(latencyInSeconds, l.Seconds())
	}
	metricName := fmt.Sprintf("Latency of GetObject(File Size: %s)", getObjectSize(object_level))
	// Making a json object, the format is something like:
	// { "metricName" : dataPoints []float64 }
	return map[string][]float64{metricName: latencyInSeconds}
}

func (s3P S3Performancer) Init() {
}

// Default performancer

type DefaultPerformancer struct {
}

func (d DefaultPerformancer) Start(ctx context.Context, params EventParams) map[string][]float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	totalVirtualMemoryInMB := "Sys(MB)"
	results := map[string][]float64{
		totalVirtualMemoryInMB: []float64{},
	}
	results[totalVirtualMemoryInMB] = append(results[totalVirtualMemoryInMB], float64(m.TotalAlloc/1024/1024))

	return results
}

func (d DefaultPerformancer) Init() {
}
