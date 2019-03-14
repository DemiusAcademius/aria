package main

// https://kubernetes.io/docs/tasks/access-application-cluster/communicate-containers-same-pod-shared-volume/

import (
	"fmt"
	"log"
	"net"
	// "os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"demius/image-builder/internal"
	"demius/image-builder/internal/api"

	// git "gopkg.in/src-d/go-git.v4"
	// "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func main() {
	config := internal.LoadConfiguration()

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", config.ServerPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	if config.CertFile == "" || config.KeyFile == "" {
		log.Fatalln("No key & crt files are specified in the environment")
	}

	creds, err := credentials.NewServerTLSFromFile(config.CertFile, config.KeyFile)
	if err != nil {
		log.Fatalf("Failed to generate credentials %v", err)
	}
	opts = []grpc.ServerOption{grpc.Creds(creds)}

	grpcServer := grpc.NewServer(opts...)
	api.RegisterImageBuilderServer(grpcServer, internal.NewServer())

	log.Printf("Starting image-builder at `localhost:%s`\n", config.ServerPort)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	/*
	_, err := git.PlainClone(config.GitLocalFolder, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: config.GitUsername,
			Password: config.GitPassword,
		},		
		URL:      config.GitProvider + config.GitRepo,
		Progress: os.Stdout,
	})
	if err != nil {
		log.Fatalf("Failed clone repository: %s", err)
	}
	*/
}
