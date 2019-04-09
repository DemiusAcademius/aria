package utils

import (
	"bytes"
	"fmt"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

// Kustomization of k8s manifests
type Kustomization struct {
	Tier     string  `yaml:"tier"`
	Ns       string  `yaml:"ns"`
	Name     string  `yaml:"name"`
	Kind     string  `yaml:"kind"`
	Service  Service `yaml:"service"`
	Schedule string  `yaml:"schedule"`
}

// Service details for proxy-manager
// See proxy-manager project, annotation aria.io/proxy-config
type Service struct {
	Listener string `yaml:"listener"`
	Default  bool   `yaml:"default"`
}

// ParseKustomization service annotation aria.io/proxy-config
func ParseKustomization(content []byte) (*Kustomization, error) {
	c := Kustomization{}
	if err := yaml.Unmarshal(content, &c); err != nil {
		return nil, fmt.Errorf("Can not unmarshar ProxyConfig: %v", err)
	}
	return &c, nil
}

type cronJobData struct {
	Ns       string
	Name     string
	Schedule string
}

// KustomizeCronJob generate cronjob manifest for k8s
func KustomizeCronJob(kustomization *Kustomization, tmpl *template.Template) ([]byte, error) {
	data := cronJobData{
		Ns:       kustomization.Ns,
		Name:     kustomization.Name,
		Schedule: kustomization.Schedule,
	}

	manifestBuffer := new(bytes.Buffer)
	err := tmpl.Execute(manifestBuffer, data)
	if err != nil {
		return nil, fmt.Errorf("can not apply variables to cronjob template: %v", err)
	}
	return manifestBuffer.Bytes(), nil
}

type deploymentData struct {
	Ns   string
	Name string
}

// KustomizeDeployment generate cronjob manifest for k8s
func KustomizeDeployment(kustomization *Kustomization, tmpl *template.Template) ([]byte, error) {
	data := deploymentData{
		Ns:   kustomization.Ns,
		Name: kustomization.Name,
	}

	manifestBuffer := new(bytes.Buffer)
	err := tmpl.Execute(manifestBuffer, data)
	if err != nil {
		return nil, fmt.Errorf("can not apply variables to deployment template: %v", err)
	}
	return manifestBuffer.Bytes(), nil
}

type serviceData struct {
	Ns       string
	Name     string
	Listener string
	Default  bool
}

// KustomizeService generate cronjob manifest for k8s
func KustomizeService(kustomization *Kustomization, tmpl *template.Template) ([]byte, error) {
	data := serviceData{
		Ns:   kustomization.Ns,
		Name: kustomization.Name,
		Listener: kustomization.Service.Listener,
		Default: kustomization.Service.Default,
	}

	manifestBuffer := new(bytes.Buffer)
	err := tmpl.Execute(manifestBuffer, data)
	if err != nil {
		return nil, fmt.Errorf("can not apply variables to service template: %v", err)
	}
	return manifestBuffer.Bytes(), nil
}

