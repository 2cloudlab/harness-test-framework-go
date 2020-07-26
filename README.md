# Test Harness Framework in Go based on AWS Lambda Function

![](test-harness-framework-go.png)

## Usage

1. Build from source

```bash
go build worker-handler.go shared-data-struct.go S3Performancer.go
```

2. Zip the generated executable bin

```bash
zip worker-handler.zip worker-handler
```

3. Provision Infrustructure

```bash
terraform init
terraform plan
terraform apply
```

4. Launch Test Harness & collect result

```bash
go run auto-run.go
```