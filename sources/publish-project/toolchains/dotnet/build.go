package dotnet

import (
	"io"
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"text/template"

	"github.com/fatih/color"

	"demius/publish-project/core"
)

// Project is header of .csproj xml file
type Project struct {
	XMLName    xml.Name        `xml:"Project"`
	Properties []PropertyGroup `xml:"PropertyGroup"`
}

// PropertyGroup  of .csproj xml file
type PropertyGroup struct {
	XMLName         xml.Name `xml:"PropertyGroup"`
	TargetFramework string   `xml:"TargetFramework"`
}

// DockerTemplate variables for Dockerfile template
type DockerTemplate struct {
	Version    string
	Executable string
}

// Build dot.net project and fill the grpc Request
func Build(configPath, projectPath, projectName string) []byte {
	projectFile := path.Join(projectPath, projectName+".csproj")
	project := loadProject(projectFile)

	targetFramework := project.Properties[0].TargetFramework
	runtimeVersion := targetFramework[len(targetFramework)-3:]

	core.PrintBlue("target framework: ", targetFramework)
	core.PrintBlue(" runtime version: ", runtimeVersion)

	println()
	color.Magenta("DOTNET PUBLISH")

	cmd := exec.Command("dotnet", "publish", "-c", "Release")
	cmd.Dir = projectPath

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%s\n", out.String())
		core.PrintErrorAndPanic(err)
	}

	publishPath := path.Join(projectPath, "bin", "Release", targetFramework, "publish")
	println()
	core.PrintBlue("    publish path: ", publishPath)

	assetsPath := path.Join(projectPath, "assets")
	if core.FileExistsAndDir(assetsPath) {
		core.PrintBlue("     assets path: ", assetsPath)
		if err = copyAssets(assetsPath, path.Join(publishPath, "assets")); err != nil {
			core.PrintErrorAndPanic(err)
		}
	}

	println()
	color.Magenta("GENERATE TARBALL")

	dockerfile, customDockerfileExist := loadCustomDockerfile(projectPath)
	if !customDockerfileExist {
		dockerfile = generateDockerfile(configPath, projectName, runtimeVersion)
	}

	tarBuffer, err := core.CreateTarball(publishPath, dockerfile)
	if err != nil {
		core.PrintErrorAndPanic(err)
	}
	return tarBuffer
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

func loadCustomDockerfile(projectPath string) ([]byte, bool) {
	dockerfilePath := path.Join(projectPath, "Dockerfile")
	if core.FileExists(dockerfilePath) {
		fp, err := os.Open(dockerfilePath)
		if err != nil {
			core.PrintErrorAndPanic(fmt.Errorf("can not open source file %s: %v", dockerfilePath, err))
		}
		defer fp.Close()
	
		dockerfile, err := ioutil.ReadAll(fp)
		if err != nil {
			core.PrintErrorAndPanic(fmt.Errorf("can not read Dockerfile %s: %v", dockerfilePath, err))
		}
		return dockerfile, true
	}
	return nil, false
}

func generateDockerfile(configPath, projectName, runtimeVersion string) []byte {
	dockerfilePath := path.Join(configPath, "dockerfiles", "dotnet", "Dockerfile")
	fp, err := os.Open(dockerfilePath)
	if err != nil {
		core.PrintErrorAndPanic(fmt.Errorf("can not open source file %s: %v", dockerfilePath, err))
	}
	defer fp.Close()

	dockerfile, err := ioutil.ReadAll(fp)
	if err != nil {
		core.PrintErrorAndPanic(fmt.Errorf("can not read Dockerfile %s: %v", dockerfilePath, err))
	}

	dt := DockerTemplate{Version: runtimeVersion, Executable: projectName}
	tmpl, err := template.New("Dockerfile").Parse(string(dockerfile[:len(dockerfile)]))
	if err != nil {
		core.PrintErrorAndPanic(fmt.Errorf("can not parse template for Dockerfile: %v", err))
	}
	dockerfileBuffer := new(bytes.Buffer)
	err = tmpl.Execute(dockerfileBuffer, dt)
	if err != nil {
		core.PrintErrorAndPanic(fmt.Errorf("can not apply variables to Dockerfile template: %v", err))
	}
	return dockerfileBuffer.Bytes()
}

func copyAssets(sourcePath, destinationPath string) error {
	sourceLen := len(sourcePath) + 1

	if _, err := os.Stat(destinationPath); os.IsNotExist(err) {
		os.Mkdir(destinationPath, os.ModePerm)
	}

	return filepath.Walk(sourcePath, func(filePath string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}

        if !f.Mode().IsRegular() {
            return nil
		}
		
		source, err := os.Open(filePath)
        if err != nil {
			return fmt.Errorf("Can not open source file %s: %v", filePath, err)
        }
        defer source.Close()		

		relationalPath := filePath[sourceLen:]

		fullDestinationPath := path.Join(destinationPath, relationalPath)
		destination, err := os.Create(fullDestinationPath)
        if err != nil {
			return fmt.Errorf("Can not create destination file %s: %v", fullDestinationPath, err)
        }
        defer destination.Close()
		if _, err = io.Copy(destination, source); err != nil {
			return fmt.Errorf("Can not copy from source to destination: %v", err)
		}

		return nil
	})
}
