package publisher

import (
	"strconv"
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
	imageName := calcImageName(request)
	response := &api.Response{ ImageName: imageName }	

	apiNamespaces := s.Clientset.CoreV1().Namespaces()

	if _, err := apiNamespaces.Get(request.Namespace, kubemeta.GetOptions{}); err != nil {
		return resp_error(response, fmt.Sprintf("Could not get namespace `%s`, got error '%s'\n", request.Namespace, err.Error())), nil
	}

	newVersion := time.Now().Format("0601021504") // stdYear stdZeroMonth stdZeroDay stdHour stdZeroMinute
	nvi,_ := strconv.Atoi(newVersion)
	newVersion = strings.ToUpper( fmt.Sprintf(strconv.FormatInt(int64(nvi), 16)) )

	newImageName := s.RegistryURL + "/" + imageName + ":" + newVersion

	switch request.Kind {
	case api.ArtifactKind_CronJob:
		{
			batchAPI := s.Clientset.BatchV1beta1()
			apiJobs := batchAPI.CronJobs(request.Namespace)

			job, err := apiJobs.Get(request.Name, kubemeta.GetOptions{})
			if err != nil {
				return resp_error(response, fmt.Sprintf("Could not get job `%s`, got error '%s'\n", request.Name, err.Error())), nil
			}
			containers := job.Spec.JobTemplate.Spec.Template.Spec.Containers
			idx := searchContainerIdxByImageName(containers, imageName)
			if idx == -1 {
				return resp_error(response, fmt.Sprintf("Could not get container with image `%s`\n", imageName)), nil
			}

			if err = PushImage(newImageName, request.DockerContent, s.Credentials); err != nil {
				return resp_error(response, err.Error()), nil
			}

			job.Spec.JobTemplate.Spec.Template.Spec.Containers[idx].Image = newImageName
			if _, err := apiJobs.Update(job); err != nil {
				return resp_error(response, fmt.Sprintf("job update error `%v`\n", err)), nil
			}
		}
	case api.ArtifactKind_Deployment:
		{
			appsAPI := s.Clientset.AppsV1()
			apiDeployments := appsAPI.Deployments(request.Namespace)

			deployment, err := apiDeployments.Get(request.Name, kubemeta.GetOptions{})
			if err != nil {
				return resp_error(response, fmt.Sprintf("Could not get deployment `%s`, got error '%s'\n", request.Name, err.Error())), nil
			}
			containers := deployment.Spec.Template.Spec.Containers
			idx := searchContainerIdxByImageName(containers, imageName)
			if idx == -1 {
				return resp_error(response, fmt.Sprintf("Could not get container with image `%s`\n", imageName)), nil
			}

			if err = PushImage(newImageName, request.DockerContent, s.Credentials); err != nil {
				return resp_error(response, err.Error()), nil
			}

			deployment.Spec.Template.Spec.Containers[idx].Image = newImageName
			if _, err := apiDeployments.Update(deployment); err != nil {
				return resp_error(response, fmt.Sprintf("deployment update error `%v`\n", err)), nil
			}
		}
	}

	return resp_ok(response, "0x" + newVersion), nil
}

func resp_error(response *api.Response, errorDesc string) *api.Response {
	response.ResponseVariants = &api.Response_ErrorDescription{
		ErrorDescription: errorDesc,
	}
	return response
}

func resp_ok(response *api.Response, version string) *api.Response {
	response.ResponseVariants = &api.Response_ImageVersion{
		ImageVersion: version,
	}
	return response
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

func calcImageName(request *api.Request) string {
	imageName := request.Name + "."
	if request.Kind == api.ArtifactKind_Deployment {
		imageName += request.Tier + "."
	}
	imageName += request.Namespace
	return imageName
}