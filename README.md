# Test Harness Framework in Go based on AWS Lambda Function

![](test-harness-framework-go.png)

## Introduction

This is a Test Harness Framework(written in Go) based on AWS Lambda Function. It can be used in the following scenarios:

1. Launch a large number of loaders to do performance tests against your services.
2. Do performance tests on AWS services, such as S3, DynamoDB etc.

All you require to do are:

1. Write a code snippet for your scenario in Go.
2. Tune the test parameters & start the tests.

After that, a few reports, including stats and raw reports, will be automatically generated in `.csv` format. You can import it into sheet to compare the benchmarks.

## Usage

1. Build from source

```bash
cd test-harness-framework-go
dep ensure
GOOS="linux" GOARCH="amd64" go build -ldflags "-w -s" worker-handler.go shared-data-struct.go *Performancer.go
```

2. Zip the generated executable bin

```bash
zip worker-handler.zip worker-handler
```

3. Provision Infrustructure

```bash
terraform init
terraform apply -var="bucket_name=<replace-with-your-bucket-name>"
```

4. Launch Test Harness & collect reports

```bash
go run auto-run.go shared-data-struct.go -bucket-name <your-provisioned-bucket-name-in-step-3>
```

5. Destroy resources

```bash
terraform destroy -var="bucket_name=<replace-with-your-bucket-name>"
```