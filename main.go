package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"

	"github.com/drone/drone-go/drone"
)

type CloudwatchClient interface {
	PutMetricData(ctx context.Context, params *cloudwatch.PutMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.PutMetricDataOutput, error)
}

// Ensure required env vars are set
// This doesn't scale well, but we only have a few variables
// Makes it much easier to debug Lambda failures and has the added bonus
// of failing much faster
func verifyEnvVars() error {
	for _, v := range []string{
		"DRONE_TOKEN",
		"DRONE_SERVER",
		"CLOUDWATCH_METRICS_NAMESPACE",
	} {
		if os.Getenv(v) == "" {
			fmt.Printf("Required env var '%s' not set\n", v)
			return errors.New("Missing env var")
		}
	}
	return nil
}

// Retrieve all builds, pending or running
// We'll filter downstream
func getQueuedBuilds(c drone.Client) []*drone.Stage {
	s, err := c.Queue()

	if err != nil {
		fmt.Printf("Error retrieving build queue %s\n", err.Error())
		os.Exit(1)
	}

	return s
}

func reportBuilds(c drone.Client, cw CloudwatchClient, builds []*drone.Stage) {

	// If there aren't any pending builds, exit
	// The cloudwatch alarm needs to treat missing data as not breaching
	if len(builds) < 1 {
		fmt.Println("Build queue is empty")
		os.Exit(0)
	}

	// Iterate through pending builds
	for _, b := range builds {

		// Running builds are good
		// A running build doesn't need a new worker node
		if b.Status != "running" {

			// Create dimensions array for each queued build
			var dimensions []types.Dimension

			// b.Labels is a map[string]string representing the builds node labels
			for k, v := range b.Labels {
				// Build CW metric dimensions using build's node labels
				dimensions = append(dimensions, types.Dimension{Name: aws.String(k), Value: aws.String(v)})

			}

			// Write metric for this queued build to Cloudwatch
			putCloudwatchMetric(
				cw,
				dimensions,
			)

		}
	}

}

func putCloudwatchMetric(c CloudwatchClient, d []types.Dimension) error {

	md := []types.MetricDatum{
		{
			MetricName:        aws.String("QueuedBuilds"),
			Dimensions:        d,
			Value:             aws.Float64(1.0),
			StorageResolution: aws.Int32(60),
			Unit:              types.StandardUnitCount,
		},
	}

	p := cloudwatch.PutMetricDataInput{
		Namespace:  aws.String(os.Getenv("CLOUDWATCH_METRICS_NAMESPACE")),
		MetricData: md,
	}

	_, err := c.PutMetricData(context.TODO(), &p)

	if err != nil {
		fmt.Printf("Error putting metric data - %s\n", err.Error())
		return err
	} else {
		fmt.Println("PutMetric success!")
		return nil
	}
}

func handler(ctx context.Context, e events.CloudWatchEvent) {
	// Verify env vars so we fail fast if any are missing
	if err := verifyEnvVars(); err != nil {
		os.Exit(1)
	}

	// Create clients
	droneClient := newDroneClient()
	cwClient := newCloudwatchClient()

	// Get queued builds
	builds := getQueuedBuilds(droneClient)

	// Report metrics to cloudwatch
	reportBuilds(droneClient, cwClient, builds)
}

func main() {
	lambda.Start(handler)
}
