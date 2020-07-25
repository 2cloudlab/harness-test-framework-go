package main

import (
	"fmt"
	"testing"
)

func TestDownloadByPrefix(t *testing.T) {
	init_shared_resource()
	results := downloadByPrefix("2cloudlab-performance-benchmark-bucket", "d5e17827-9682-41ae-a30b-1551dc4bea66")
	fmt.Println(len(results), "hjh")
}
