package dotnet

import (
	"fmt"
	"bytes"
	"os/exec"
	"log"
	"io/ioutil"
	"encoding/xml"
	"path"

	"github.com/fatih/color"

	"demius/publish-project/api"
	"demius/publish-project/core"
)

// Project is header of .csproj xml file
type Project struct {
	XMLName    xml.Name        `xml:"Project"`
	Properties []PropertyGroup `xml:"PropertyGroup"`
}

// PropertyGroup  of .csproj xml file
type PropertyGroup struct {
	XMLName xml.Name `xml:"PropertyGroup"`
	TargetFramework string  `xml:"TargetFramework"`
}

// Build dot.net project and fill the grpc Request
func Build(projectPath, projectName string, request *api.Request) {
	projectFile := path.Join(projectPath, projectName+".csproj")
	project := loadProject(projectFile)

	targetFramework := project.Properties[0].TargetFramework
	runtimeVersion := targetFramework[len(targetFramework)-3:]

	core.PrintBlue("target framework: " , targetFramework)
	core.PrintBlue(" runtime version: " , runtimeVersion)

	println()
	color.Magenta("DOTNET PUBLISH")

	cmd := exec.Command("dotnet", "publish", "-c", "Release")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%s\n", out.String())	
		core.PrintErrorAndPanic(err)
	}

	publishPath := path.Join(projectPath, "bin", "Release", targetFramework, "publish")
	println()
	core.PrintBlue("    publish path: " , publishPath)
}

func loadProject(path string) *Project {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("%s get err #%v", path, err)
	}
	var c = &Project{}

	if err = xml.Unmarshal(file, c); err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c

}
