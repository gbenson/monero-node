resource "aws_secretsmanager_secret" "graphite_secret" {
  name = "graphite_api"
}

resource "aws_iam_policy" "graphite_secret_read" {
  name = "graphite_secret_read"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
	Effect = "Allow",
	Action = "secretsmanager:GetSecretValue",
	Resource = aws_secretsmanager_secret.graphite_secret.arn
      },
    ],
  })
}
