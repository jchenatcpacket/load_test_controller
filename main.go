package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

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

	container_name := fmt.Sprintf("load_test_%d", duplicate_id)

	resp, err := cli.ContainerCreate(
		ctx,
		container_config,
		nil,
		nil,
		nil,
		container_name,
	)
	if err != nil {
		panic(err)
	}

	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
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

	var wg sync.WaitGroup
	for i := 0; i <= 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			spawnLoadTest(ctx, cli, username, password, i)
		}()
	}
	wg.Wait()
}
