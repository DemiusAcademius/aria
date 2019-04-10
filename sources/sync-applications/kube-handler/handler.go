package kubehandler

import (
	"bytes"
	"flag"
	"fmt"
	"path/filepath"

	// "k8s.io/apimachinery/pkg/api/errors"
	"github.com/fatih/color"

	apibatch "k8s.io/api/batch/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sYaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Connect to kubernetes API outside cluster
func Connect(homePath string) *kubernetes.Clientset {
	kubeconfig := flag.String("kubeconfig", filepath.Join(homePath, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}

// ManifestType enum of manifest file types
type ManifestType = int

const (
	// CronJobManifest contains k8s manifest for cronjob resource
	CronJobManifest ManifestType = iota
	// DeploymentManifest contains k8s manifest for deployment resource
	DeploymentManifest
	// ServiceManifest contains k8s manifest for service resource
	ServiceManifest
)

// ApplyManifest decode manifest from byte buffer and apply it to k8s
func ApplyManifest(clientset *kubernetes.Clientset, manifest []byte, manifestType ManifestType) error {
	decoder := k8sYaml.NewYAMLOrJSONDecoder(bytes.NewReader(manifest), 1000)

	switch manifestType {
	case CronJobManifest:
		return applyCronjob(clientset, decoder)
	case DeploymentManifest:		
		return applyDeployment(clientset, decoder)
	case ServiceManifest:
		return applyService(clientset, decoder)
	}

	return nil
}

func applyCronjob(clientset *kubernetes.Clientset, decoder *k8sYaml.YAMLOrJSONDecoder) error {
	j := &apibatch.CronJob{}

	if err := decoder.Decode(&j); err != nil {
		return err
	}
		
	decodedArtifact("      cronjob: ", j.Name, j.Namespace)

	batchAPI := clientset.BatchV1beta1()
	apiJobs := batchAPI.CronJobs(j.Namespace)

	if _, err := apiJobs.Get(j.Name, metav1.GetOptions{}); err != nil {
		// create job
		if _, err := apiJobs.Create(j); err != nil {
			return fmt.Errorf("job create error '%s'", err.Error())
		}	
		created()
	} else {
		allreadyExists()
	}


	return nil
}

func applyDeployment(clientset *kubernetes.Clientset, decoder *k8sYaml.YAMLOrJSONDecoder) error {
	d := &appsv1.Deployment{}

	if err := decoder.Decode(&d); err != nil {
		return err
	}

	decodedArtifact("   deployment: ", d.Name, d.Namespace)

	appsAPI := clientset.AppsV1()
	apiDeployments := appsAPI.Deployments(d.Namespace)

	if _, err := apiDeployments.Get(d.Name, metav1.GetOptions{}); err != nil {
		// create deployment
		if _, err := apiDeployments.Create(d); err != nil {
			return fmt.Errorf("deployment create error '%s'", err.Error())
		}
		created()
	} else {
		allreadyExists()
	}

	return nil
}

func applyService(clientset *kubernetes.Clientset, decoder *k8sYaml.YAMLOrJSONDecoder) error {
	s := &apiv1.Service{}

	if err := decoder.Decode(&s); err != nil {
		return err
	}

	decodedArtifact("      service: ", s.Name, s.Namespace)
	
	api := clientset.CoreV1()
	apiServices := api.Services(s.Namespace)
	if _, err := apiServices.Get(s.Name, metav1.GetOptions{}); err != nil {
		// create deployment
		if _, err := apiServices.Create(s); err != nil {
			return fmt.Errorf("service create error '%s'", err.Error())
		}
		created()
	} else {
		allreadyExists()
	}

	return nil
}

func decodedArtifact(artifact, name, namespace string) {
	print(color.HiBlackString(artifact))
	fmt.Printf("%s", name)
	println(color.HiBlackString(".%s", namespace))

}

func allreadyExists() {
	println(color.YellowString("     allready exists"))
}

func created() {
	println(color.GreenString("     created"))
}
