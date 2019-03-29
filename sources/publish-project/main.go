package main

import (
	"path/filepath"
	"path"
	"github.com/fatih/color"

	"demius/publish-project/api"
	"demius/publish-project/core"
	"demius/publish-project/toolchains/dotnet"
)

func main() {
	executablePath := core.ExecutableDir()
	projectPath := core.WorkingDir()
	core.PrintBlue(" executable path: " , executablePath)
	core.PrintBlue("    project path: " , projectPath)

	config := core.LoadConfig(path.Join(executablePath, "aria-config.yaml"))
	artifactConfig := core.LoadArtifactConfig(path.Join(projectPath, "artifact-config.yaml"))

	core.CheckArtifactKind(artifactConfig.ArtifactKind)

	core.PrintBlue("     aria-server: " , config.AriaServer)
	core.PrintBlue("       imagename: " , artifactConfig.ImageName)

	projectName := filepath.Base(projectPath)

	core.PrintBlue("     projectname: " , projectName)

	projectType := core.DetectProjectType(projectPath, projectName)

	request := &api.Request {
		ArtifactName: artifactConfig.ApplicationName,
		ArtifactKind: artifactConfig.ArtifactKind,
		Namespace: artifactConfig.Namespace,
		ImageName: artifactConfig.ImageName,
	}

	switch projectType {
	case core.DotnetProjectType: dotnet.Build(projectPath, projectName, request)
	default: {
		color.Red("This project type not yeat realized!")
		panic("")	
	}
	}

	println()
	color.Magenta("OK")
}

