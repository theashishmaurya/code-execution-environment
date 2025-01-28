package main

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"

	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// Docker Config for the Node container
var NodeConfig container.Config = container.Config{
	Image:        "node:alpine",
	Cmd:          []string{"node", "-v"},
	ExposedPorts: nat.PortSet{"3000/tcp": struct{}{}},
	User:         "node",
	WorkingDir:   "/app",
}

func getTemplate(templatePath string) ([]fs.DirEntry, error) {

	// Reads the files in the templates folder
	template, err := os.ReadDir(templatePath)

	if err != nil {
		return nil, fmt.Errorf("failed to read template folder: %w", err)
	}

	for _, file := range template {
		fmt.Println(file.Name())
	}

	// Returns the files in the templates folder
	return template, nil

}

func CreateTarArchive(sourceDir, outputFile string) error {
	// Create the output file
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Create a tar writer
	tw := tar.NewWriter(file)
	defer tw.Close()

	// Walk through the source directory
	err = filepath.Walk(sourceDir, func(filePath string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create a tar header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		// Update the header name to be relative to the source directory
		relPath, err := filepath.Rel(sourceDir, filePath)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(relPath)

		// Write the header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// Write the file content (if it's a file)
		if !fi.IsDir() {
			data, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer data.Close()

			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func main() {
	// Create the docker client
	client, err := client.NewClientWithOpts(client.FromEnv)

	ctx := context.Background()

	if err != nil {
		panic(err)
	}

	// 0. Pull the image if not exists.

	// _, err = client.ImagePull(ctx, "node:alpine", image.PullOptions{})
	// if err != nil {
	// 	panic(fmt.Sprintf("Error pulling Node.js image: %v", err))
	// }
	// fmt.Println("Node.js image pulled successfully.")

	// 1. Create A container

	containerResp, err := client.ContainerCreate(context.Background(), &NodeConfig, &container.HostConfig{}, &network.NetworkingConfig{}, &v1.Platform{}, "my-go-container")

	if err != nil {
		panic(fmt.Sprintf("Error creating container: %v", err))
	}

	fmt.Printf("Container created successfully! ID: %s\n", containerResp.ID)
	fmt.Printf("Warnings: %v\n", containerResp.Warnings) // 2. Install Node and NPM in the container
	err = client.ContainerStart(ctx, containerResp.ID, container.StartOptions{})

	// STEP 0 - Prepare the right container and install the runtime requires. Currently supported Nodejs / Python / Java / Golang / C++
	// STEP 1 - Expose the Container Life Cycle Via API
	// STEP 2 - Expose the shell Via API
	// STEP 3 - Expose the FS via API

}
