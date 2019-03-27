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

const (
	// CronJobFile contains k8s manifest for cronjob resource
	CronJobFile = "cronjob.yaml"
	// DeploymentFile contains k8s manifest for deployment resource
	DeploymentFile = "deployment.yaml"
	// ServiceFile contains k8s manifest for service resource
	ServiceFile = "service.yaml"
)

func main() {
	executablePath := utils.ExecutableDir()
	homePath := utils.UserHomeDir()
	fmt.Println("executable path: " + executablePath)
	fmt.Println("      home path: " + homePath)

	clientset := handler.Connect(homePath)

	source := filepath.Join(homePath, "applications")

	if err := walkApplications(source, clientset); err != nil {
		color.Red(fmt.Sprintf("%s", err))
		return
	}

	color.HiCyan("OK")
}

func walkApplications(source string, clientset *kubernetes.Clientset) error {
	prefixLen := len(source) + 1

	return filepath.Walk(source, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walk dir%s: %v", path, err)
		}
		if f.IsDir() {
			return nil
		}
		
		filename := filepath.Base(path)
		if !strings.HasSuffix(filename, ".yaml") {
			return nil
		}
		if filename == CronJobFile || filename == DeploymentFile || filename == ServiceFile {
			printArtifactPath(prefixLen, path, filename)

			filedata, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			var mt handler.ManifestType
			switch filename {
			case CronJobFile: mt = handler.CronJobManifest
			case DeploymentFile: mt = handler.DeploymentManifest
			case ServiceFile: mt = handler.ServiceManifest
			default: {
				print(color.HiBlackString("      unknown resource type\n"))
				return nil
			}
					}
			if err = handler.ApplyManifest(clientset, filedata, mt); err != nil {
				return err
			}
		}

		return nil
	})
}

func printArtifactPath(prefixLen int, path, filename string) {
	strippedPath := path[prefixLen:len(path)-len(filename)-1]
	print(color.HiBlackString("sync: "))
	fmt.Printf("%s", strippedPath)
	print(color.HiBlackString(" / "))
	print(color.GreenString("%s", filename[:len(filename)-5]))
	println(color.HiBlackString(" .yaml"))
}
