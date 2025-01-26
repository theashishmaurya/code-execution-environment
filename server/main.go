package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func main() {
	// Create the docker client
	client, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		panic(err)
	}

	// List all containers
	containers, err := client.ContainerList(context.Background(), container.ListOptions{})

	if err != nil {
		panic(err)
	}

	for _, ctr := range containers {
		fmt.Printf("%s %s\n", ctr.ID, ctr.Image)
	}

}
