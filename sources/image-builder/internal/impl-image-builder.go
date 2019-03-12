package internal

import (
	"context"

	"demius/image-builder/internal/api"
)

type imageBuilderServer struct {
	initialized bool
	// applications map[string]*pb.Application
}

// NewServer create new grpc server
func NewServer() api.ImageBuilderServer {
	s := &imageBuilderServer{
		initialized: false,
		// applications: make(map[string]*pb.Application),
	}
	return s
}

var okResponse = api.Response{Code: api.ResponseCode_OK}

// Init connect to Git repo and clone it
func (s *imageBuilderServer) Init(ctx context.Context, gitRepo *api.GitRepo) (*api.Response, error) {
	s.initialized = true

	// TODO: save git repo params and clone

	return &okResponse, nil
}

// Pull fetch Git repo and merge it
func (s *imageBuilderServer) Pull(ctx context.Context, _ *api.Empty) (*api.Response, error) {

	// TODO: save git repo params and clone

	return &okResponse, nil
}

// BuildProject from cloned git repo
func (s *imageBuilderServer) BuildProject(ctx context.Context, project *api.Project) (*api.BuildResponse, error) {

	// TODO: build project from cloned git repo

	response := api.BuildResponse {
		Code: api.ResponseCode_OK,
		Project: project,
	}

	return &response, nil
}

// BuildNamespace from cloned git repo for all projects in namespace and returns response for each project
func (s *imageBuilderServer) BuildNamespace(namespace *api.Namespace, stream api.ImageBuilder_BuildNamespaceServer) error {

	// TODO: build project from cloned git repo

	return nil
}

// BuildAll from cloned git repo for all projects in all namespaces and returns response for each project
func (s *imageBuilderServer) BuildAll(_ *api.Empty, stream api.ImageBuilder_BuildAllServer) error {

	// TODO: build project from cloned git repo

	return nil
}
