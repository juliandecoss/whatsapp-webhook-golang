variable "app_version" {}

variable "s3_bucket" {
  default = "backend-lambda-golang"
}

variable "retention_in_days"{
  default = 1
}

variable "aws_access_key" {
  description = "Access Key for AWS"
}

variable "aws_secret_key" {
  description = "Secret Key for AWS"
}

variable "whatsapp_token" {
  description = "permanent token to send whatsapps"
}