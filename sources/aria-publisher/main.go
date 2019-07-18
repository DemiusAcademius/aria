package main

// https://kubernetes.io/docs/tasks/access-application-cluster/communicate-containers-same-pod-shared-volume/

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"demius/aria-publisher/api"
	"demius/aria-publisher/publisher"
)

// MaxMessageSize maximum message size of GRPC
const MaxMessageSize = 1024 * 1024 * 24

func main() {
	config := publisher.LoadConfiguration()

	registryAuthPath := "/auth/credentials.yaml"
	regCreds := publisher.LoadRegistryCredentials(registryAuthPath)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", config.ServerPort))
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
	opts = []grpc.ServerOption{grpc.Creds(creds), grpc.MaxRecvMsgSize(MaxMessageSize)}

	grpcServer := grpc.NewServer(opts...)
	api.RegisterPublishRequestServer(grpcServer, publisher.NewServer(regCreds, config.RegistryURL))

	log.Printf("Starting aria-publish-service at `localhost:%s`\n", config.ServerPort)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
