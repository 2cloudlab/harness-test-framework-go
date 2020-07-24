package main

import (
	"context"
	"fmt"
)

type S3Performancer struct {
}

func (s3P S3Performancer) Start(ctx context.Context, params EventParams) []byte {
	for i := 0; i < params.CountInSingleInstance; i++ {
		fmt.Println("Hello World!")
	}
	return []byte("")
}

func (s3P S3Performancer) Init() {
}
