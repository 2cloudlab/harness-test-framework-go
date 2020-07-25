GOOS="linux" GOARCH="amd64" go build -ldflags "-w -s" test-harness-framework.go shared-data-struct.go
zip test-harness-framework.zip test-harness-framework

GOOS="linux" GOARCH="amd64" go build -ldflags "-w -s" worker-handler.go shared-data-struct.go *Performancer.go
zip worker-handler.zip worker-handler

aws-vault exec slz -- go run auto-run.go shared-data-struct.go 2cloudlab-performance-benchmark-bucket