variable "admin_ssh_key" {
  description = "SSH public key for server admin"
  type        = string
}

variable "aws_lambda_python_powertools_layer_arn" {
  description = "arn of AWS Lambda Powertools (Python) layer to use"
  type        = string
  default     = "arn:aws:lambda:us-east-1:017000801446:layer:AWSLambdaPowertoolsPythonV2:46"
}
