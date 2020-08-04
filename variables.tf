variable "worker_handler_zip_file" {
  description = ""
  type        = string
  default     = "worker-handler"
}

variable "bucket_name" {
  description = ""
  type        = string
}

variable "worker_runtime" {
  description = ""
  type        = string
  default     = "go1.x"
}
