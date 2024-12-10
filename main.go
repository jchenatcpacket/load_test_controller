package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/joho/godotenv"
)

func spawnLoadTest(ctx context.Context, cli *client.Client, username string, password string, duplicate_id int) {
	load_test_image := "load_test"

	container_config := &container.Config{
		Image: load_test_image,
		Cmd:   []string{"opam", "exec", "dune", "exec", "./client_piaf.exe"},
		Tty:   false,
		Env: []string{
			"OCAMLRUNPARAM=b",
			"REQUEST_NUMBER=10000",
			"URL=https://10.51.42.143/api/epg_fr/monitors/",
			fmt.Sprintf("USERNAME=%s", username),
			fmt.Sprintf("PASSWORD=%s", password),
		},
	}

	resp, err := cli.ContainerCreate(
		ctx,
		container_config,
		nil,
		nil,
		nil,
		fmt.Sprintf("load_test_%d", duplicate_id),
	)

	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	for i := 0; i <= 2; i++ {
		spawnLoadTest(ctx, cli, username, password, i)
	}
}
