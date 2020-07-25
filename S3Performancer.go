package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// S3 performancer

type S3Performancer struct {
}

type dataPoint struct {
	TotalSizeInBytes int
	Latency          float64
}

func (s3P S3Performancer) Start(ctx context.Context, params EventParams) []byte {
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

func (s3P S3Performancer) Init() {
}

// default performancer

type DefaultPerformancer struct {
}

func (s3P DefaultPerformancer) Start(ctx context.Context, params EventParams) []byte {
	return []byte("")
}

func (s3P DefaultPerformancer) Init() {
}
