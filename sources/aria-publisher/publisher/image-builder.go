package publisher

import (
	"bytes"
	"os"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// PushImage guild docker image and push it to registry
func PushImage(imageName string, image []byte, credentials *RegistryCredentials) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return fmt.Errorf("Could not create docker client, got error '%s'", err.Error())
	}

	imageContext := bytes.NewReader(image)
	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		imageContext,
		types.ImageBuildOptions{
			Dockerfile: "Dockerfile",
			Tags:       []string{imageName},
			Remove:     true,
		})
	if err != nil {
		return fmt.Errorf("Could not build image, got error '%s'", err.Error())
	}

	defer imageBuildResponse.Body.Close()
	lines, err := ioutil.ReadAll(imageBuildResponse.Body)
	if err != nil {
		return fmt.Errorf("Could not read build image response, got error '%s'", err.Error())
	}
	os.Stdout.Write(lines)

	authConfig := types.AuthConfig{
		Username: credentials.UserName,
		Password: credentials.Password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return fmt.Errorf("Could not marshal registry credentials, got error '%s'", err.Error())
	}

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	imagePushResponse, err := cli.ImagePush(ctx, imageName, types.ImagePushOptions{
		RegistryAuth: authStr,
	})
	if err != nil {
		return fmt.Errorf("Could not push image to registry, got error '%s'", err.Error())
	}

	defer imagePushResponse.Close()
	if _, err = ioutil.ReadAll(imagePushResponse); err != nil {
		return fmt.Errorf("Could not read push image response, got error '%s'", err.Error())
	}

	return nil
}
