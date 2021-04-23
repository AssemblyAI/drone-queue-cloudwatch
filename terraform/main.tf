locals {
  common_tags = {
    owner   = "engineering/ops"
    Project = "drone"
    tier    = "monitoring"
  }
}

provider "aws" {
  region = "us-west-2"
}


terraform {
  required_version = "~> 0.14.0"

  required_providers {
    aws = {
      version = "3.37.0"
      source  = "hashicorp/aws"
    }
  }
}