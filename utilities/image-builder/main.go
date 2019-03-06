package main

import (
	"log"
	"os"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func main() {
	var config = LoadConfiguration()
	println("Starting image-builder")

	_, err := git.PlainClone(config.GitLocalFolder, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: config.GitUsername,
			Password: config.GitPassword,
		},		
		URL:      config.GitProvider + config.GitRepo,
		Progress: os.Stdout,
	})
	if err != nil {
		log.Fatalf("Failed clone repository: %s", err)
	}

	println("ok")
}
