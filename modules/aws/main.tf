resource "aws_s3_bucket" "lambda_bucket" {
  bucket = "xmrig-status-listener"
}

resource "aws_s3_bucket_ownership_controls" "lambda_bucket" {
  bucket = aws_s3_bucket.lambda_bucket.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "lambda_bucket" {
  depends_on = [aws_s3_bucket_ownership_controls.lambda_bucket]

  bucket = aws_s3_bucket.lambda_bucket.id
  acl    = "private"
}

data "archive_file" "lambda_listener_zip" {
  type = "zip"

  source_file = "listener/lambda_function.py"
  output_path = ".build/listener-lambda.zip"
}

resource "aws_s3_object" "lambda_listener_zip" {
  bucket = aws_s3_bucket.lambda_bucket.id

  key    = "status-listener-lambda.zip"
  source = data.archive_file.lambda_listener_zip.output_path

  etag = filemd5(data.archive_file.lambda_listener_zip.output_path)
}

resource "aws_lambda_function" "status_listener" {
  function_name = "xmrig-status-listener"

  s3_bucket = aws_s3_bucket.lambda_bucket.id
  s3_key    = aws_s3_object.lambda_listener_zip.key

  runtime = "python3.9"
  layers = [
    var.aws_lambda_python_powertools_layer_arn,
    "arn:aws:lambda:us-east-1:919648353655:layer:python39-requests2_31_0:2",
  ]
  handler = "lambda_function.lambda_handler"

  source_code_hash = data.archive_file.lambda_listener_zip.output_base64sha256

  role = aws_iam_role.lambda_exec.arn
}

resource "aws_cloudwatch_log_group" "status_listener" {
  name = "/aws/lambda/${aws_lambda_function.status_listener.function_name}"

  retention_in_days = 30
}

resource "aws_iam_role" "lambda_exec" {
  name = "XMRigStatusListenerLambda"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Sid    = ""
      Principal = {
	Service = "lambda.amazonaws.com"
      }
    }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_policy" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_graphite_secret_read" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = aws_iam_policy.graphite_secret_read.arn
}

resource "aws_apigatewayv2_api" "lambda" {
  name          = "xmrig-status-listener"
  protocol_type = "HTTP"
  description   = "XMRig status listener"
}

resource "aws_apigatewayv2_stage" "lambda" {
  api_id = aws_apigatewayv2_api.lambda.id

  name        = "v1"
  auto_deploy = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gw.arn

    format = jsonencode({
      requestId               = "$context.requestId"
      sourceIp                = "$context.identity.sourceIp"
      requestTime             = "$context.requestTime"
      protocol                = "$context.protocol"
      httpMethod              = "$context.httpMethod"
      resourcePath            = "$context.resourcePath"
      routeKey                = "$context.routeKey"
      status                  = "$context.status"
      responseLength          = "$context.responseLength"
      integrationErrorMessage = "$context.integrationErrorMessage"
      }
    )
  }
}

resource "aws_apigatewayv2_integration" "status_listener" {
  api_id = aws_apigatewayv2_api.lambda.id

  integration_uri    = aws_lambda_function.status_listener.invoke_arn
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
}

resource "aws_apigatewayv2_route" "status_listener" {
  api_id = aws_apigatewayv2_api.lambda.id

  route_key = "POST /recv"
  target    = "integrations/${aws_apigatewayv2_integration.status_listener.id}"
}

resource "aws_cloudwatch_log_group" "api_gw" {
  name = "/aws/api_gw/${aws_apigatewayv2_api.lambda.name}"

  retention_in_days = 30
}

resource "aws_lambda_permission" "api_gw" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.status_listener.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.lambda.execution_arn}/*/*"
}
