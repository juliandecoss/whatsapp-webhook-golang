resource "aws_lambda_function" "whatsapp-webhook" {
  function_name = "whatsapp-webhook-wedding"

  # The bucket name as created earlier with "aws s3api create-bucket"
  s3_bucket = var.s3_bucket
  s3_key    = "v${var.app_version}/${var.file_zip_name}.zip"

  # "main" is the filename within the zip file (main.js) and "handler"
  # is the name of the property under which the handler function was
  # exported in that file.
  handler = "bootstrap"
  runtime = "provided.al2"

  role = aws_iam_role.lambda_exec.arn

  environment {
    variables = {
      WHATSAPP_TOKEN = var.whatsapp_token
    }
  }

}

# IAM role which dictates what other AWS services the Lambda function
# may access.
resource "aws_iam_role" "lambda_exec" {
  name = "lambda-policy"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF

  # Política adicional para permitir escritura en CloudWatch Logs
  inline_policy {
    name = "lambda-cloudwatch-logs-policy"

    policy = jsonencode({
      Version = "2012-10-17",
      Statement = [
        {
          Action = [
            "logs:CreateLogGroup",
            "logs:CreateLogStream",
            "logs:PutLogEvents",
            "dynamodb:PutItem"
          ],
          Effect   = "Allow",
          Resource = "*"
        }
      ]
    })
  }
}

resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.whatsapp-webhook.function_name
  principal     = "apigateway.amazonaws.com"

  # The /*/* portion grants access from any method on any resource
  # within the API Gateway "REST API".
  source_arn = "${aws_api_gateway_rest_api.webhook.execution_arn}/*/*"
}

resource "aws_cloudwatch_log_group" "log_group" {
  name              = "/aws/lambda/${aws_lambda_function.whatsapp-webhook.function_name}"
  retention_in_days = var.retention_in_days
}
