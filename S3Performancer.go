package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3 performancer
type S3Performancer struct {
}

func (s3P S3Performancer) Start(ctx context.Context, params EventParams) []byte {
	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(params.RawJson), &m)
	if err != nil {
		recordError(err)
		return []byte{}
	}
	object_level := int(m["FileSize"].(float64))
	sample_data_key := getObjectName(object_level)
	testTasks := make(chan int, params.CountInSingleInstance)
	samples := make(chan int, params.CountInSingleInstance)
	for g := 0; g < params.CountInSingleInstance; g++ {
		go func(tasks <-chan int, results chan<- int) {
			for range tasks {
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

				results <- 1
			}
		}(testTasks, samples)
	}

	benchmarkTimer := time.Now()

	for i := 0; i < params.CountInSingleInstance; i++ {
		testTasks <- i
	}

	close(testTasks)

	for s := 0; s < params.CountInSingleInstance; s++ {
		_ = <-samples
	}

	totalObjectSizeInBytes := 1024 * (1 << (object_level - 1)) * params.CountInSingleInstance
	// stop the timer for this benchmark
	totalTime := time.Now().Sub(benchmarkTimer)

	return []byte(fmt.Sprintf(`[{"TotalObjectSizeInBytes": %d, "TotalTime": %f}]`, totalObjectSizeInBytes, totalTime.Seconds()))
}

func (s3P S3Performancer) Init() {
}

// Default performancer

type DefaultPerformancer struct {
}

type dataPoint struct {
	TotalSizeInBytes int
	Latency          float64
}

func (d DefaultPerformancer) Start(ctx context.Context, params EventParams) []byte {
	results := []dataPoint{}
	for i := 0; i < params.CountInSingleInstance; i++ {
		rand.Seed(time.Now().UnixNano())
		min := 10
		max := 30
		sizeInBytes := rand.Intn(max-min+1) + min
		fmt.Println(sizeInBytes)
		laytency := rand.Float64()
		fmt.Println(laytency)
		rand.Seed(time.Now().UnixNano())
		singleData := dataPoint{
			TotalSizeInBytes: sizeInBytes,
			Latency:          laytency,
		}
		results = append(results, singleData)
	}

	b, _ := json.Marshal(results)
	return b
}

func (d DefaultPerformancer) Init() {
}
