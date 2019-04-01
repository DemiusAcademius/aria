package publisher

import (
	"log"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

// RegistryCredentials docker-registry info
type RegistryCredentials struct {
	UserName string `yaml:"username"`
    Password string `yaml:"password"`
}

// LoadRegistryCredentials from yaml file
func LoadRegistryCredentials(registryAuthPath string) *RegistryCredentials {    
    yamlFile, err := ioutil.ReadFile(registryAuthPath)
    if err != nil {
        log.Printf("%s Get err   #%v ",registryAuthPath, err)
    }
    c := &RegistryCredentials {}

    err = yaml.Unmarshal(yamlFile, c)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return c
}