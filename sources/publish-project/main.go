package main

import (
	"os"
	"time"
	"context"
	"fmt"
	"path/filepath"
	"path"

	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"demius/publish-project/api"
	"demius/publish-project/core"
	"demius/publish-project/toolchains/dotnet"
	"demius/publish-project/toolchains/ui"
)

// MaxMessageSize maximum message size of GRPC
const MaxMessageSize = 1024 * 1024 * 12

func main() {
	configPath := path.Join(core.UserHomeDir(),"PublishProject")
	core.PrintBlue("     config path: " , configPath)
	config := core.LoadConfig(path.Join(configPath, "aria-config.yaml"))
	core.PrintBlue("     aria-server: " , config.AriaServer)

	filepath.Walk(core.WorkingDir(), func(projectPath string, f os.FileInfo, err error) error {
		if err != nil {
			core.PrintErrorAndPanic(fmt.Errorf("error walk dir %s: %v", projectPath, err))
		}

		if !f.IsDir() {
			return nil
		}

		isProjectFolder, artifactConfigPath := detectProjectFolder(projectPath)
		if isProjectFolder {
			publishProject(config, configPath, projectPath, artifactConfigPath)
		}

		return nil
	})
}

func detectProjectFolder(projectPath string) (bool,string) {
	artifactConfigPath := path.Join(projectPath, "artifact-config.yaml")
	if core.FileExists(artifactConfigPath) {
        return true, artifactConfigPath
    }
	return false, ""
}

func publishProject(config *core.Config, configPath, projectPath, artifactConfigPath string) {
	println()
	core.PrintBlue("    project path: " , projectPath)

	artifactConfig := core.LoadArtifactConfig(artifactConfigPath)

	artifactKind := core.ConvertArtifactKind(artifactConfig.ArtifactKind)

	core.PrintBlue("            tier: " , artifactConfig.Tier)

	projectName := filepath.Base(projectPath)

	core.PrintBlue("     projectname: " , projectName)

	projectType := core.DetectProjectType(projectPath, projectName)

	request := &api.Request {
		Name: artifactConfig.ApplicationName,
		Kind: artifactKind,
		Namespace: artifactConfig.Namespace,
		Tier: artifactConfig.Tier,
	}

	switch projectType {
	case core.DotnetProjectType:
		request.DockerContent = dotnet.Build(configPath, projectPath, projectName)
	case core.WebUIProjectType:
		request.DockerContent = ui.Build(configPath, projectPath)
	default:
		color.Red("     This project type not yeat realized!")
		return
	}
	
	println()
	core.PrintBlueExtended("      image size: ", fmt.Sprintf("%v",len(request.DockerContent)), " bytes")

	uploadToServer(configPath, config.AriaServer, request)

	println()
	color.Magenta("OK")
}

func uploadToServer(configPath, ariaServer string, request *api.Request) {
	println()
	color.Magenta("UPLOAD TO SERVER")

	creds, err := credentials.NewClientTLSFromFile(path.Join(configPath,"acc.io.crt"), "")
	if err != nil {
		core.PrintErrorAndPanic(fmt.Errorf("Failed to create TLS credentials %v", err))
	}
	grpcCredentials := grpc.WithTransportCredentials(creds)

	conn, err := grpc.Dial(
			ariaServer, 
			grpcCredentials,
			grpc.WithMaxMsgSize(MaxMessageSize),
			grpc.WithTimeout(60*time.Second),
		)
	if err != nil {
		core.PrintErrorAndPanic(fmt.Errorf("Failed to dial applications-server %v", err))
	}
	defer conn.Close()
	client := api.NewPublishRequestClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	response, err := client.Publish(ctx, request)
	if err != nil {
		core.PrintErrorAndPanic(fmt.Errorf("Upload to server error: %v", err))
	}

	println()
	core.PrintBlue("      image name: ", response.ImageName)

	errorDescription := response.GetErrorDescription()
	if errorDescription != "" {
		core.PrintInRedAndPanic("Error in applications-server: " + errorDescription)
	}

	core.PrintBlue("   image version: ", response.GetImageVersion())
}

