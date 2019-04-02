package ui

import (
	"os"
	"fmt"
	"bytes"
	"os/exec"
	"io/ioutil"
	"path"

	"github.com/fatih/color"

	"demius/publish-project/core"
)

// Build react.js project with yarn and fill the grpc Request
func Build(configPath, projectPath string) []byte {
	println()
	color.Magenta("YARN BUILD")

	cmd := exec.Command("yarn", "build")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%s\n", out.String())	
		core.PrintErrorAndPanic(err)
	}

	publishPath := path.Join(projectPath, "build")
	println()
	core.PrintBlue("      build path: " , publishPath)

	println()
	color.Magenta("GENERATE TARBALL")

	dockerfile := generateDockerfile(configPath)
	tarBuffer, err := core.CreateTarball(publishPath, dockerfile)
	if err != nil {
		core.PrintErrorAndPanic(err)
	}
	return tarBuffer
}

func generateDockerfile(configPath string) []byte {
	dockerfilePath := path.Join(configPath, "dockerfiles", "ui", "Dockerfile")
	fp, err := os.Open(dockerfilePath)
	if err != nil {
		core.PrintErrorAndPanic(fmt.Errorf("can not open source file %s: %v", dockerfilePath, err))
	}
	defer fp.Close()

	dockerfile, err := ioutil.ReadAll(fp)
	if err != nil {
		core.PrintErrorAndPanic(fmt.Errorf("can not read Dockerfile %s: %v", dockerfilePath, err))
	}

	return dockerfile
}
