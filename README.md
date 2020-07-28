# Test Harness Framework in Go based on AWS Lambda Function

![](test-harness-framework-go.png)

This is a Test Harness Framework(written in Go) based on AWS Lambda Function. It can be used in the following scenarios:

1. Launch a large number of loaders to do performance tests against your services.
2. Do performance tests on AWS services, such as S3, DynamoDB etc.

All you require to do are:

1. Write a code snippet for your scenario in Go.
2. Tune the test parameters & start the tests.

After that, a few reports, such as `raw-data-<TaskName>-<DateTime>-<TaskId>.csv` and `report-<TaskName>-<DateTime>-<TaskId>.csv`, will be automatically generated in `reports` folder. `report-<TaskName>-<DateTime>-<TaskId>.csv` file contains some stats information, such as avg, min, max, p25, p50, p75, p90 and p99, which are calculated beyond the `raw-data-<TaskName>-<DateTime>-<TaskId>.csv` file. In addition, it will merged reports base on the same `TaskName` but with different test conditions, the merged reports name is something like `raw-data-<TaskName>-<DateTime>.csv` and `report-<TaskName>-<DateTime>.csv`, you can import them into sheet to compare the benchmarks or visualize them.

## 1. Prerequisites

Before you use the framework, please install the following tools on top of **Linux OS**, and pay attention to their versions.

* Install Go, and make sure its version is at least `v1.14`. After successfully install, type the following command to make sure you are on the right version:

```bash
$ go version
go version go1.14.2 linux/amd64
```

* Install Terraform, and make sure its version is `v0.12.19`. Type the following command to make sure you are on the right version:

```bash
terraform version
Terraform v0.12.19
```

* Install [dep](https://golang.github.io/dep/docs/installation.html)

* Clone the framework and prepare dependencies

Only execute the following command at the first time when using Test Harness Framework.

```bash
git clone https://github.com/2cloudlab/test-harness-framework-go.git
cd test-harness-framework-go
dep ensure
```

## 2. Write your `*Performancer.go`



## 3. Usage

* Build from source

```bash
make build
```

* Provision Infrustructure

```bash
make auto_provision BUCKET_NAME="<replace-with-your-bucket-name>"
```

* Launch Test Harness & collect reports

The following command will start your Tasks in parallel, you should tell it from where(`BUCKET_NAME="test-reports-repository"`) and when(`TIME_TO_WAIT="2"`) to start collecting reports.

```bash
make run BUCKET_NAME="<replace-with-your-bucket-name>" TIME_TO_WAIT="<time-to-wait-before-collecting-reports-in-minute>"
```

* Destroy resources

If you no longer use the provisioned resources in step 2, make sure to call the following command to destroy them, so that you are not charged by AWS.

```bash
make auto_destroy BUCKET_NAME="<replace-with-your-bucket-name>"
```