package core

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// Config for connect to aria-publisher
type Config struct {
	AriaServer string `yaml:"aria-server"`
}

// ArtifactConfig connect local project with k8s artifact
type ArtifactConfig struct {
	Namespace       string `yaml:"namespace"`
	ApplicationName string `yaml:"app-name"`
	ArtifactKind    string `yaml:"kind"`
	Tier            string `yaml:"tier"`
}

// LoadConfig from specified filepath
func LoadConfig(path string) *Config {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("%s get err #%v", path, err)
	}
	var c = &Config{}

	if err = yaml.Unmarshal(file, c); err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

// LoadArtifactConfig from specified filepath
func LoadArtifactConfig (path string) *ArtifactConfig {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("%s get err #%v", path, err)
	}
	var c = &ArtifactConfig{}

	if err = yaml.Unmarshal(file, c); err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

