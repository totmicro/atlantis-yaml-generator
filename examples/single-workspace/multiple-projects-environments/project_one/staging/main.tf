# Define the provider block to specify the cloud provider you're using
provider "aws" {
  region = "us-east-1"
}

# Create a simple AWS S3 bucket
resource "aws_s3_bucket" "example_bucket" {
  bucket = "my-terraform-hello-world-bucket"
}

# Output the bucket name
output "bucket_name" {
  value = aws_s3_bucket.example_bucket.id
}
