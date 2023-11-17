resource "aws_s3_bucket" "xmrig_listener" {
  bucket = "xmrig-status-listener"
}

resource "aws_s3_bucket_ownership_controls" "xmrig_listener" {
  bucket = aws_s3_bucket.xmrig_listener.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "xmrig_listener" {
  depends_on = [aws_s3_bucket_ownership_controls.xmrig_listener]

  bucket = aws_s3_bucket.xmrig_listener.id
  acl    = "private"
}

data "archive_file" "listener_lambda_zip" {
  type = "zip"

  source_dir  = "listener"
  output_path = "${path.module}/listener-lambda.zip"
}

resource "aws_s3_object" "listener_lambda_zip" {
  bucket = aws_s3_bucket.xmrig_listener.id

  key    = "listener-lambda.zip"
  source = data.archive_file.listener_lambda_zip.output_path

  etag = filemd5(data.archive_file.listener_lambda_zip.output_path)
}
