variable "TF_VAR_aws_region" {}
variable "TF_VAR_aws_access_key" {}
variable "TF_VAR_aws_secret_key" {}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 3.0, < 4.0"
    }
  }
}

provider "aws" {
  region     = var.TF_VAR_aws_region
  access_key = var.TF_VAR_aws_access_key
  secret_key = var.TF_VAR_aws_secret_key
}

resource "aws_s3_bucket_object" "valid_elements" {
  bucket = "city-validation-s3"
  key    = "results/valid_elements.json"
  acl    = "private"
  content = file("${path.module}/results/valid_elements.json")
}

resource "aws_s3_bucket_object" "invalid_elements" {
  bucket = "city-validation-s3"
  key    = "results/invalid_elements.json"
  acl    = "private"
  content = file("${path.module}/results/invalid_elements.json")
}

resource "aws_s3_bucket_object" "unprocessable_files" {
  bucket = "city-validation-s3"
  key    = "results/unprocessable_files.csv"
  acl    = "private"
  content = file("${path.module}/results/unprocessable_files.csv")
}
