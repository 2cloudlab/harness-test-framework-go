
terraform {
  required_version = "= 0.12.19"
}

provider "aws" {
  version = "= 2.58"
  # region  = "ap-northeast-1"
}

data "aws_caller_identity" "current" {}

data "aws_region" "current" {}

locals {
  account_id                 = "${data.aws_caller_identity.current.account_id}"
  region_name                = "${data.aws_region.current.name}"
  time_out_in_second         = 900
  test_harness_function_name = "test-harness-framework"
  _worker_handler_config_object = jsondecode(file("worker-handler-config.json"))
  worker_handlers_function_name = {
    for i in range(local._worker_handler_config_object["MinFunctionMemoryInMB"],local._worker_handler_config_object["MaxFunctionMemoryInMB"]+1,local._worker_handler_config_object["IncreaseMemoryByInMB"]):
    format("%d", i) => format("%s-%d", local._worker_handler_config_object["FunctionNamePrefix"], i)
  }
  policy_for_test_harness    = <<POLICY
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "logs:CreateLogGroup",
            "Resource": "arn:aws:logs:${local.region_name}:${local.account_id}:*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogStream",
                "logs:PutLogEvents"
            ],
            "Resource": [
                "arn:aws:logs:${local.region_name}:${local.account_id}:log-group:/aws/lambda/${local.test_harness_function_name}:*"
            ]
        },
        {
            "Action": [
                "lambda:InvokeFunction"
            ],
            "Resource": "arn:aws:lambda:${local.region_name}:${local.account_id}:function:${var.worker_handler_zip_file}*",
            "Effect": "Allow"
        }
    ]
}
POLICY
  policy_for_worker_handler  = <<POLICY
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "logs:CreateLogGroup",
            "Resource": "arn:aws:logs:${local.region_name}:${local.account_id}:*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogStream",
                "logs:PutLogEvents"
            ],
            "Resource": [
                "arn:aws:logs:${local.region_name}:${local.account_id}:log-group:/aws/lambda/${var.worker_handler_zip_file}*:*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "s3:PutObject",
                "s3:GetObject"
            ],
            "Resource": "arn:aws:s3:::${var.bucket_name}/*"
        }
    ]
}
POLICY
}

# create IAM Role for test harness framework Lambda Function

resource "aws_iam_role" "role_for_test_harness" {
  name = "test_harness_role"

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
POLICY
}

resource "aws_iam_policy" "permission_policy_for_test_harness" {
  name        = "permission_policy_for_test_harness"
  description = "Permissions policy for test harness"

  policy = local.policy_for_test_harness
}

resource "aws_iam_role_policy_attachment" "policy_role_attach" {
  role       = aws_iam_role.role_for_test_harness.name
  policy_arn = aws_iam_policy.permission_policy_for_test_harness.arn
}

# create test harness framework Lambda Function

resource "aws_lambda_function" "test_harness_lambda" {
  filename      = "${local.test_harness_function_name}.zip"
  function_name = local.test_harness_function_name
  role          = aws_iam_role.role_for_test_harness.arn
  handler       = local.test_harness_function_name
  timeout       = local.time_out_in_second
  memory_size   = 128

  # The filebase64sha256() function is available in Terraform 0.11.12 and later
  # For Terraform 0.11.11 and earlier, use the base64sha256() function and the file() function:
  # source_code_hash = "${base64sha256(file("lambda_function_payload.zip"))}"
  source_code_hash = filebase64sha256("${local.test_harness_function_name}.zip")

  # check out the detail runtime from here https://docs.aws.amazon.com/lambda/latest/dg/API_CreateFunction.html#SSS-CreateFunction-request-Runtime
  runtime = "go1.x"
}

# create IAM Role for worker handler function

resource "aws_iam_role" "role_for_worker_handler" {
  name = "worker_handler_role"

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
POLICY
}

resource "aws_iam_policy" "permission_policy_for_worker_handler" {
  name        = "permission_policy_for_worker_handler"
  description = "Permissions policy for worker handler"

  policy = local.policy_for_worker_handler
}

resource "aws_iam_role_policy_attachment" "policy_role_attach_for_worker_handler" {
  role       = aws_iam_role.role_for_worker_handler.name
  policy_arn = aws_iam_policy.permission_policy_for_worker_handler.arn
}

# create worker handler Lambda Function

resource "aws_lambda_function" "worker_handler_lambda" {
  for_each = local.worker_handlers_function_name
  filename      = "${var.worker_handler_zip_file}.zip"
  function_name = each.value
  role          = aws_iam_role.role_for_worker_handler.arn
  handler       = var.worker_handler_zip_file
  timeout       = local.time_out_in_second
  memory_size   = tonumber(each.key)

  # The filebase64sha256() function is available in Terraform 0.11.12 and later
  # For Terraform 0.11.11 and earlier, use the base64sha256() function and the file() function:
  # source_code_hash = "${base64sha256(file("lambda_function_payload.zip"))}"
  source_code_hash = filebase64sha256("${var.worker_handler_zip_file}.zip")

  environment {
    variables = {
      BUCKET_NAME = var.bucket_name
    }
  }

  # check out the detail runtime from here https://docs.aws.amazon.com/lambda/latest/dg/API_CreateFunction.html#SSS-CreateFunction-request-Runtime
  runtime = var.worker_runtime
}

resource "aws_s3_bucket" "b" {
  bucket = var.bucket_name
  acl    = "private"
  force_destroy = true
}

resource "aws_s3_bucket_object" "object" {
  for_each     = fileset("test-data/", "*")
  bucket = aws_s3_bucket.b.id
  key    = "test-data/${each.value}"
  source = "test-data/${each.value}"

  # The filemd5() function is available in Terraform 0.11.12 and later
  # For Terraform 0.11.11 and earlier, use the md5() function and the file() function:
  # etag = "${md5(file("path/to/file"))}"
  etag = filemd5("test-data/${each.value}")
}