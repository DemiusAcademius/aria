package main

import (
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
	projectPath := core.WorkingDir()
	configPath := path.Join(core.UserHomeDir(),"PublishProject")
	core.PrintBlue("     config path: " , configPath)
	core.PrintBlue("    project path: " , projectPath)

	config := core.LoadConfig(path.Join(configPath, "aria-config.yaml"))
	artifactConfig := core.LoadArtifactConfig(path.Join(projectPath, "artifact-config.yaml"))

	artifactKind := core.ConvertArtifactKind(artifactConfig.ArtifactKind)

	core.PrintBlue("     aria-server: " , config.AriaServer)
	core.PrintBlue("       imagename: " , artifactConfig.ImageName)

	projectName := filepath.Base(projectPath)

	core.PrintBlue("     projectname: " , projectName)

	projectType := core.DetectProjectType(projectPath, projectName)

	request := &api.Request {
		Name: artifactConfig.ApplicationName,
		Kind: artifactKind,
		Namespace: artifactConfig.Namespace,
		ImageName: artifactConfig.ImageName,
	}

	switch projectType {
	case core.DotnetProjectType:
		request.DockerContent = dotnet.Build(configPath, projectPath, projectName)
	case core.WebUIProjectType:
		request.DockerContent = ui.Build(configPath, projectPath)
	default:
		core.PrintInRedAndPanic("This project type not yeat realized!")
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
	errorDescription := response.GetErrorDescription()
	if errorDescription != "" {
		core.PrintInRedAndPanic("Error in applications-server: " + errorDescription)
	}

	println()
	core.PrintBlue("   image version: ", response.GetImageVersion())
}

