package ui

import (
	"strings"
	"bufio"
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
	cmd.Dir = projectPath
	
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

	buildVersion := extractVersion(path.Join(projectPath, "build.gradle"))
	core.PrintBlue("   build version: " , buildVersion)

	println()
	color.Magenta("GENERATE TARBALL")

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

func extractVersion(gradlefilePath string) string {
	fp, err := os.Open(gradlefilePath)
	if err != nil {
		core.PrintErrorAndPanic(fmt.Errorf("can not open file %s: %v", gradlefilePath, err))
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line,"version") {
			versionString := strings.Split(line," ")[1]
			return versionString[1:len(versionString)-2]
		}
    }

    if err := scanner.Err(); err != nil {
        core.PrintErrorAndPanic(fmt.Errorf("can not read file %s: %v", gradlefilePath, err))
	}
	
	return "1.0"
}