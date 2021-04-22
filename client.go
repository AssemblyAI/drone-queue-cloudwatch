package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/drone/drone-go/drone"
	"golang.org/x/oauth2"
)

func newDroneClient() drone.Client {
	conf := new(oauth2.Config)
	a := conf.Client(oauth2.NoContext, &oauth2.Token{AccessToken: os.Getenv("DRONE_TOKEN")})

	c := drone.NewClient(os.Getenv("DRONE_SERVER"), a)

	return c
}

func newCloudwatchClient() *cloudwatch.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		fmt.Printf("Failed to load AWS configuration, %v", err)
		os.Exit(1)
	}

	return cloudwatch.NewFromConfig(cfg)
}
