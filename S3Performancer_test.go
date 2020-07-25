package main

import (
	"fmt"
	"testing"
	"time"
)

type MockContext struct{}

func (mctx MockContext) Deadline() (deadline time.Time, ok bool) {
	return time.Now(), true
}

func (mctx MockContext) Done() <-chan struct{} {
	return nil
}

func (mctx MockContext) Err() error {
	return nil
}

func (mctx MockContext) Value(key interface{}) interface{} {
	return nil
}

func TestS3PerformancerStart(t *testing.T) {
	init_shared_resource()
	performer := S3Performancer{}
	params := EventParams{CountInSingleInstance: 2, RawJson: `{ "FileSize" : 1}`}
	context := MockContext{}
	fmt.Println(string(performer.Start(context, params)))
}

func TestS3PerformancerInit(t *testing.T) {
	init_shared_resource()
	performer := S3Performancer{}
	performer.Init()
}
