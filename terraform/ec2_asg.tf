// "Standard" ASG
resource "aws_autoscaling_group" "drone_standard_ec2_asg" {
  name                      = "drone-standard-workers"
  max_size                  = 4
  min_size                  = 2
  desired_capacity          = 2
  vpc_zone_identifier       = ["some subnet id", "some other subnet id"]
  default_cooldown          = 300
  health_check_grace_period = 120
  health_check_type         = "EC2"
  termination_policies      = ["OldestLaunchTemplate", "ClosestToNextInstanceHour", "Default"]
  suspended_processes       = null

  launch_template {
    id      = "some launch template id"
    version = "$Latest"
  }

  instance_refresh {
    strategy = "Rolling"
    preferences {
      instance_warmup        = 120
      min_healthy_percentage = 50
    }
  }

}

resource "aws_autoscaling_policy" "standard_drone_queue_up" {
  name                   = "drone-build-queue-depth-up"
  scaling_adjustment     = 2 // Launch 2
  adjustment_type        = "ChangeInCapacity"
  policy_type            = "SimpleScaling"
  cooldown               = 180 // 3 minutes
  autoscaling_group_name = aws_autoscaling_group.drone_standard_ec2_asg.name
}

resource "aws_autoscaling_policy" "standard_drone_queue_down" {
  name                   = "drone-build-queue-depth-down"
  scaling_adjustment     = -1 // Remove 1
  adjustment_type        = "ChangeInCapacity"
  policy_type            = "SimpleScaling"
  cooldown               = 900 // 15 minutes
  autoscaling_group_name = aws_autoscaling_group.drone_standard_ec2_asg.name
}


resource "aws_cloudwatch_metric_alarm" "standard_worker_queue_depth" {
  alarm_name                = "drone-standard-build-queue-depth"
  comparison_operator       = "GreaterThanOrEqualToThreshold"
  evaluation_periods        = "2"
  metric_name               = "QueuedBuilds"
  namespace                 = "Drone"
  period                    = "60"
  statistic                 = "Sum"
  threshold                 = "3"
  datapoints_to_alarm       = 2
  alarm_description         = "Drone build queue depth for standard workers"
  insufficient_data_actions = []
  treat_missing_data        = "notBreaching"
  actions_enabled           = true

  tags = {
    env           = "cicd"
    Name          = "drone-standard-build-queue-depth"
    owner         = "engineering"
    tier          = "cicd"
    "drone/type"  = "worker"
    "drone/class" = "standard"
  }

  dimensions = {
    "os"    = "linux"
    "class" = "standard"
  }

  alarm_actions = [
    aws_autoscaling_policy.standard_drone_queue_up.arn
  ]

  ok_actions = [ 
    aws_autoscaling_policy.standard_drone_queue_down.arn
  ]
}


// GPU ASG
resource "aws_autoscaling_group" "drone_gpu_ec2_asg" {
  name                      = "drone-gpu-workers"
  max_size                  = 3
  min_size                  = 1
  desired_capacity          = 1
  vpc_zone_identifier       = ["some subnet id", "some other subnet id"]
  default_cooldown          = 300
  health_check_grace_period = 120
  health_check_type         = "EC2"
  termination_policies      = ["OldestLaunchTemplate", "ClosestToNextInstanceHour", "Default"]
  suspended_processes       = null

  launch_template {
    id      = "some launch template id for GPU nodes"
    version = "$Latest"
  }

  instance_refresh {
    strategy = "Rolling"
    preferences {
      instance_warmup        = 120
      min_healthy_percentage = 50
    }
  }

}

resource "aws_autoscaling_policy" "gpu_drone_queue_up" {
  name                   = "drone-build-queue-depth-up"
  scaling_adjustment     = 2 // Launch 2
  adjustment_type        = "ChangeInCapacity"
  policy_type            = "SimpleScaling"
  cooldown               = 180 // 3 minutes
  autoscaling_group_name = aws_autoscaling_group.drone_gpu_ec2_asg.name
}

resource "aws_autoscaling_policy" "gpu_drone_queue_down" {
  name                   = "drone-build-queue-depth-down"
  scaling_adjustment     = -1 // Remove 1
  adjustment_type        = "ChangeInCapacity"
  policy_type            = "SimpleScaling"
  cooldown               = 900 // 15 minutes
  autoscaling_group_name = aws_autoscaling_group.drone_gpu_ec2_asg.name
}

resource "aws_cloudwatch_metric_alarm" "gpu_worker_queue_depth" {
  alarm_name                = "drone-gpu-build-queue-depth"
  comparison_operator       = "GreaterThanOrEqualToThreshold"
  evaluation_periods        = "2"
  metric_name               = "QueuedBuilds"
  namespace                 = "Drone"
  period                    = "60"
  statistic                 = "Sum"
  threshold                 = "3"
  datapoints_to_alarm       = 2
  alarm_description         = "Drone build queue depth for GPU workers"
  insufficient_data_actions = []
  treat_missing_data        = "notBreaching"
  actions_enabled           = true

  tags = {
    env           = "cicd"
    Name          = "drone-gpu-build-queue-depth"
    owner         = "engineering"
    tier          = "cicd"
    "drone/type"  = "worker"
    "drone/class" = "gpu"
  }

  dimensions = {
    "os"    = "linux"
    "class" = "gpu"
  }

  alarm_actions = [
    aws_autoscaling_policy.standard_drone_queue_up.arn
  ]

  ok_actions = [ 
    aws_autoscaling_policy.standard_drone_queue_down.arn
  ]
}