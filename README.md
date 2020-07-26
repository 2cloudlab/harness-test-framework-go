# Test Harness Framework in Go based on AWS Lambda Function

![](test-harness-framework-go.png)

## Usage

1. Build from source

```bash
GOOS="linux" GOARCH="amd64" go build -ldflags "-w -s" worker-handler.go shared-data-struct.go *Performancer.go
```

2. Zip the generated executable bin

```bash
zip worker-handler.zip worker-handler
```

3. Provision Infrustructure

```bash
terraform init
terraform plan
terraform apply -var="bucket_name=<replace-with-your-bucket-name>"
```

4. Launch Test Harness & collect reports

```bash
go run auto-run.go shared-data-struct.go -bucket-name <your-provisioned-bucket-name-in-step-3>
```