package main

import (
	"context"
	"fmt"
	"dagger.io/dagger"
)


func main() {
	if !checkFileExists("../docker/Dockerfile.dockerfile") {
		fmt.Println("Error: Dockerfile is missing.")
		return
	}

	// Create a shared context
	ctx := context.Background()

	// Run the stages of the pipeline
	if err := Build(ctx); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
}

func Build(ctx context.Context) error {
	// Initialize Dagger client
	client, err := dagger.Connect(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	// Reads the Dockerfile and builds the container
	container := client.Host().Directory("..").DockerBuild()

	// Runs a check inside the built container to verify FastAPI loads
	_, err = container.WithExec([]string{
		"python", "-c", "import src.deployment.fastapi",
	}).Sync(ctx)

	if err != nil {
		return err
	}

	fmt.Println("Dagger succesfully built and verified the project container.")
	return nil
}
