package publisher

import (
	"context"
	"demius/aria-publisher/api"
)

type publisherServer struct {
}

// NewServer create new grpc server
func NewServer() api.PublishRequestServer {
	s := &publisherServer{
	}
	return s
}

// Publish construct docker image and update k8s artifact (deployment or cronjob)
func (s *publisherServer) Publish(ctx context.Context, request *api.Request) (*api.Response, error) {
	return &api.Response{
		ResponseVariants: &api.Response_ErrorDescription {
			ErrorDescription: "Not yet realized",
		},
	}, nil
}

