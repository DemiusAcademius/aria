package ui

import (
	"text/template"
	"os"
	"fmt"
	"bytes"
	"os/exec"
	"io/ioutil"
	"path"

	"github.com/fatih/color"

	"demius/publish-project/core"
)

type DockerTemplate struct {
	Version string
	Executable    string
}

// Build Java project with gradle and fill the grpc Request
func Build(configPath, projectPath string, projectName string) []byte {
	println()
	color.Magenta("GRADLE BUILD")

	cmd := exec.Command("gradle", "build")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%s\n", out.String())	
		core.PrintErrorAndPanic(err)
	}

	publishPath := path.Join(projectPath, "build", "libs")
	println()
	core.PrintBlue("      build path: " , publishPath)

	// TODO: analize build.gradle, read version

	println()
	color.Magenta("GENERATE TARBALL")

	// TODO: add build version to Dockerfile params
	buildVersion := "1.1"

	dockerfile := generateDockerfile(configPath, projectName, buildVersion)
	tarBuffer, err := core.CreateTarball(publishPath, dockerfile)
	if err != nil {
		core.PrintErrorAndPanic(err)
	}
	return tarBuffer
}

func generateDockerfile(configPath, projectName, buildVersion string) []byte {
	dockerfilePath := path.Join(configPath, "dockerfiles", "java", "Dockerfile")
	fp, err := os.Open(dockerfilePath)
	if err != nil {
		core.PrintErrorAndPanic(fmt.Errorf("can not open source file %s: %v", dockerfilePath, err))
	}
	defer fp.Close()

	dockerfile, err := ioutil.ReadAll(fp)
	if err != nil {
		core.PrintErrorAndPanic(fmt.Errorf("can not read Dockerfile %s: %v", dockerfilePath, err))
	}

	dt := DockerTemplate { Version: buildVersion, Executable: projectName }
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
