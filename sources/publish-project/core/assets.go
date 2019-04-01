package core

import (
	"os"
	"github.com/fatih/color"
    "path"    

    "demius/publish-project/api"
)

// ProjectType of artifact for build
type ProjectType int

const (
    // DotnetProjectType for dot.net core
    DotnetProjectType ProjectType = iota
    // WebUIProjectType for web-ui react
    WebUIProjectType
    // JavaProjectType for java gradle
    JavaProjectType
)

// DetectProjectType analize project dir and detect it type
func DetectProjectType(projectPath, projectName string) ProjectType {
    if FileExists(path.Join(projectPath, projectName + ".csproj")) {
        return DotnetProjectType
    }
    if FileExists(path.Join(projectPath, "package.json")) {
        return WebUIProjectType
    }
    if FileExists(path.Join(projectPath,"build.gradle")) {
        return JavaProjectType
    }
	color.Red("Can not detect project type")
    os.Exit(-1)
    
    return -1
}

// ConvertArtifactKind and panic if it wrong
func ConvertArtifactKind(artifactKind string) api.ArtifactKind{
    switch artifactKind {
    case "cronjob": return api.ArtifactKind_CronJob
    case "deployment": return api.ArtifactKind_Deployment
    }
    
    color.Red("Invalid artifact kind. Must be `deployment` or `cronjob`")
    os.Exit(-1)    

    return -1
}