variable "worker_handler_zip_file" {
  description = ""
  type        = string
  default     = "worker-handler"
}

variable "bucket_name" {
  description = ""
  type = string
  default = "2cloudlab-performance-benchmark-bucket"
}

variable "memory_size_in_MB" {
  description = ""
  type = number
  default = 128
}