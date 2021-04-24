# drone-queue-cloudwatch

## Overview

`drone-queue-cloudwatch` is a Lambda function that publishes queued Drone builds to Cloudwatch Metrics. Metric dimensions are built from the node labels specified by a given build.

This tool allows you to leverage AWS autoscaling across multiple Drone worker groups. Some of our CI workloads need to run on GPU instances, others don't.

## Usage

### Drone

In order for the code to actually publish metrics, you must use [node routing](https://docs.drone.io/pipeline/docker/syntax/routing/).

When `drone-queue-cloudwatch` inspects builds, it uses the node labels to build the Cloudwatch metrics dimensions.

### Lambda

This is meant to be run as a cron job.

See the `terraform/` directory for reference.

The zip artifact is accessible in the `aai-oss` S3 bucket in us-west-2
- To use a specific version of the code, use the `drone-queue-cloudwatch/<commit sha>.zip` object key
- To use the latest version, use the `drone-queue-cloudwatch/latest.zip` object key

## Autoscaling

**Note:** You want to use the "Sum" statistic when configuring autoscaling. You should also treat missing data as "Not Breaching"

This application assumes each Drone worker group passes the same `DRONE_RUNNER_LABELS` to all workers.

For example, one group would have these labels:

- `class=standard`
- `os=linux`

And another these:

- `class=gpu`
- `os=linux`

Your Drone file would specify either

```yml
node:
  os: linux
  class: gpu
```

or 

```yml
node:
  os: linux
  class: standard
```

Your autoscaling trigger would launch more instances / containers for a given worker group based on how many queued builds there are.

