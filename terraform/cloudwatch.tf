resource "aws_cloudwatch_event_rule" "cron" {
  name                = "drone-queue-cloudwatch"
  description         = "Check pending builds and publish metrics"
  is_enabled          = true
  schedule_expression = "rate(1 minute)"

  tags = local.common_tags
}

resource "aws_cloudwatch_event_target" "cron" {
  rule = aws_cloudwatch_event_rule.cron.name
  arn  = aws_lambda_function.lambda.arn
}


resource "aws_cloudwatch_log_group" "log_group" {
  name              = "/aws/lambda/${aws_lambda_function.lambda.function_name}"
  retention_in_days = 1
  tags              = local.common_tags
}
