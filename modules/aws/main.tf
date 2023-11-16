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
