resource "aws_lambda_function" "lambda" {
  function_name                  = "drone-queue-cloudwatch"
  role                           = aws_iam_role.lambda.arn
  handler                        = "main"
  s3_bucket                      = var.artifact_bucket
  s3_key                         = var.artifact_key
  runtime                        = "go1.x"
  memory_size                    = 128
  reserved_concurrent_executions = 1
  timeout                        = 5
  package_type                   = "Zip"

  tags = local.common_tags

  environment {
    variables = {
      DRONE_TOKEN                  = var.drone_token
      DRONE_SERVER                 = var.drone_server
      CLOUDWATCH_METRICS_NAMESPACE = var.cloudwatch_metrics_namespace
    }
  }
  
  // Ignore changes to the source
  // CD will make updates to that and we don't want Terraform interfering
  lifecycle {
    ignore_changes = [
      s3_bucket,
      s3_key
    ]
  }
}

resource "aws_lambda_function_event_invoke_config" "config" {
  function_name                = aws_lambda_function.lambda.function_name
  maximum_event_age_in_seconds = 60
  maximum_retry_attempts       = 0
  qualifier                    = "$LATEST"
}