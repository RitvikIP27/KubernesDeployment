terraform {
  backend "s3" {
    bucket         = "ritvik-kant-statefile-lock-bucket"
    key            = "ritvik/terraform.tfstate"
    region         = "ap-south-1"
    encrypt        = true
    dynamodb_table = "terraform-lock"
  }
}