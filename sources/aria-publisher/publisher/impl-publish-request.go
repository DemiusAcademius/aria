package publisher

import (
	"context"
	"fmt"
	"strings"
	"time"

	kubecore "k8s.io/api/core/v1"
	kubemeta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"demius/aria-publisher/api"
)

/*
require (
	k8s.io/api kubernetes-1.13.4
	k8s.io/apimachinery kubernetes-1.13.4
	k8s.io/client-go v10.0.0
)
*/

type publisherServer struct {
	Clientset   *kubernetes.Clientset
	Credentials *RegistryCredentials
	RegistryURL string
}

// NewServer create new grpc server
func NewServer(credentials *RegistryCredentials, registryURL string) api.PublishRequestServer {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	s := &publisherServer{Clientset: clientset, Credentials: credentials, RegistryURL: registryURL}
	return s
}

// Publish construct docker image and update k8s artifact (deployment or cronjob)
func (s *publisherServer) Publish(ctx context.Context, request *api.Request) (*api.Response, error) {
	apiNamespaces := s.Clientset.CoreV1().Namespaces()

	if _, err := apiNamespaces.Get(request.Namespace, kubemeta.GetOptions{}); err != nil {
		return errorResponse(fmt.Sprintf("Could not get namespace `%s`, got error '%s'\n", request.Namespace, err.Error())), nil
	}

	newVersion := time.Now().Format("0601021504") // stdYear stdZeroMonth stdZeroDay stdHour stdZeroMinute
	newImageName := s.RegistryURL + "/" + request.ImageName + ":" + newVersion

	switch request.Kind {
	case api.ArtifactKind_CronJob:
		{
			batchAPI := s.Clientset.BatchV1beta1()
			apiJobs := batchAPI.CronJobs(request.Namespace)

			job, err := apiJobs.Get(request.Name, kubemeta.GetOptions{})
			if err != nil {
				return errorResponse(fmt.Sprintf("Could not get job `%s`, got error '%s'\n", request.Name, err.Error())), nil
			}
			containers := job.Spec.JobTemplate.Spec.Template.Spec.Containers
			idx := searchContainerIdxByImageName(containers, request.ImageName)
			if idx == -1 {
				return errorResponse(fmt.Sprintf("Could not get container with image `%s`\n", request.ImageName)), nil
			}

			if err = PushImage(newImageName, request.DockerContent, s.Credentials); err != nil {
				return errorResponse(err.Error()), nil
			}

			job.Spec.JobTemplate.Spec.Template.Spec.Containers[idx].Image = newImageName
			if _, err := apiJobs.Update(job); err != nil {
				return errorResponse(fmt.Sprintf("job update error `%v`\n", err)), nil
			}
		}
	case api.ArtifactKind_Deployment:
		{
			appsAPI := s.Clientset.AppsV1()
			apiDeployments := appsAPI.Deployments(request.Namespace)

			deployment, err := apiDeployments.Get(request.Name, kubemeta.GetOptions{})
			if err != nil {
				return errorResponse(fmt.Sprintf("Could not get deployment `%s`, got error '%s'\n", request.Name, err.Error())), nil
			}
			containers := deployment.Spec.Template.Spec.Containers
			idx := searchContainerIdxByImageName(containers, request.ImageName)
			if idx == -1 {
				return errorResponse(fmt.Sprintf("Could not get container with image `%s`\n", request.ImageName)), nil
			}

			if err = PushImage(newImageName, request.DockerContent, s.Credentials); err != nil {
				return errorResponse(err.Error()), nil
			}

			deployment.Spec.Template.Spec.Containers[idx].Image = newImageName
			if _, err := apiDeployments.Update(deployment); err != nil {
				return errorResponse(fmt.Sprintf("deployment update error `%v`\n", err)), nil
			}
		}
	}

	return versionResponse(newVersion), nil
}

func errorResponse(errorDesc string) *api.Response {
	return &api.Response{
		ResponseVariants: &api.Response_ErrorDescription{
			ErrorDescription: errorDesc,
		},
	}
}

func versionResponse(version string) *api.Response {
	return &api.Response{
		ResponseVariants: &api.Response_ImageVersion{
			ImageVersion: version,
		},
	}
}

func searchContainerIdxByImageName(containers []kubecore.Container, imageName string) int {
	for i, container := range containers {
		containerImageChunks := strings.Split(container.Image, "/")

		var imageNameWithVersion string
		if len(containerImageChunks) == 1 {
			imageNameWithVersion = containerImageChunks[0]
		} else {
			imageNameWithVersion = containerImageChunks[1]
		}

		containerImageName := strings.Split(imageNameWithVersion, ":")[0]
		if containerImageName == imageName {
			return i
		}
	}
	return -1
}
