package main

import (
	"io/ioutil"
	"strings"
	"path/filepath"
	"os"
	"fmt"

	"github.com/fatih/color"

	"k8s.io/client-go/kubernetes"

	handler "demius/sync-applications/kube-handler"
	utils "demius/sync-applications/utils"
)

/*
Indicate the matching versions of `k8s.io/client-go`, `k8s.io/api`, and `k8s.io/apimachinery` your project requires:

```sh
go get k8s.io/client-go@v10.0.0              # replace v10.0.0 with the required version (or use kubernetes-1.x.y tags if desired)
go get k8s.io/api@kubernetes-1.13.4          # replace kubernetes-1.13.4 with the required version
go get k8s.io/apimachinery@kubernetes-1.13.4 # replace kubernetes-1.13.4 with the required version

*/

func main() {
	executablePath := utils.ExecutableDir()
	homePath := utils.UserHomeDir()
	fmt.Println("executable path: " + executablePath)
	fmt.Println("      home path: " + homePath)

	clientset := handler.Connect(homePath)

	source := filepath.Join(homePath, "applications")
	base := filepath.Join(source, "base")

	templates, err := utils.LoadTemplates(base)
	if err != nil {
		color.Red(fmt.Sprintf("%s", err))
		return
	}

	if err := walkApplications(source, base, templates, clientset); err != nil {
		color.Red(fmt.Sprintf("%s", err))
		return
	}

	color.HiCyan("OK")
}

func walkApplications(source, base string, templates utils.Templates, clientset *kubernetes.Clientset) error {
	prefixLen := len(source) + 1

	return filepath.Walk(source, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walk dir%s: %v", path, err)
		}
		if f.IsDir() {
			return nil
		}

		if strings.HasPrefix(path, base) {
			return nil
		}
		
		filename := filepath.Base(path)
		if filename != "kustomization.yaml" {
			return nil
		}

		printArtifactPath(prefixLen, path, filename)
		filedata, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		kustomization, err := utils.ParseKustomization(filedata)
		if err != nil {
			return err
		}

		// TODO: apply kustomization to templates
		if kustomization.Kind == "cronjob" {
			tmpl := templates[utils.CronJobKind][""]
			manifest, err := utils.KustomizeCronJob(kustomization, tmpl)
			if err != nil {
				return err
			}
			if err = handler.ApplyManifest(clientset, manifest, handler.CronJobManifest); err != nil {
				return err
			}
		} else if kustomization.Kind == "deployment" {
			deploymentTmpl := templates[utils.DeploymentKind][kustomization.Tier]
			deploymentMt, err := utils.KustomizeDeployment(kustomization, deploymentTmpl)
			if err != nil {
				return err
			}
			if err = handler.ApplyManifest(clientset, deploymentMt, handler.DeploymentManifest); err != nil {
				return err
			}

			serviceTmpl := templates[utils.ServiceKind][kustomization.Tier]
			serviceMt, err := utils.KustomizeService(kustomization, serviceTmpl)
			if err != nil {
				return err
			}
			if err = handler.ApplyManifest(clientset, serviceMt, handler.ServiceManifest); err != nil {
				return err
			}
		} else {
			print(color.HiBlackString("      unknown kind\n"))
			return nil
		}

		return nil
	})
}

func printArtifactPath(prefixLen int, path, filename string) {
	strippedPath := path[prefixLen:len(path)-len(filename)-1]
	print(color.HiBlackString("sync: "))
	fmt.Printf("%s", strippedPath)
	print(color.HiBlackString(" / "))
	println(color.HiBlackString("%s", filename))
}
