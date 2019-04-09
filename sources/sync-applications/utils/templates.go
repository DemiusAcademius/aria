package utils

import (
	"io/ioutil"
	"github.com/fatih/color"
	"os"
	"strings"
	"fmt"
	"path/filepath"
	"text/template"
)

// ArtifactKind is kind of template: deployment, service, cronjob
type ArtifactKind = int

const (
	// DeploymentKind for deployment template
	DeploymentKind ArtifactKind = iota
	// ServiceKind for service template
	ServiceKind
	// CronJobKind for cronjob template
	CronJobKind
)

const (
	// CronJobName contains k8s manifest for cronjob resource
	CronJobName = "cronjob"
	// DeploymentName contains k8s manifest for deployment resource
	DeploymentName = "deployment"
	// ServiceName contains k8s manifest for service resource
	ServiceName = "service"
)

// TemplatesByTier map from tiers (ui, api etc) to templates
type TemplatesByTier = map[string]*template.Template
// Templates map for artifact kinds to templates
type Templates = map[ArtifactKind]TemplatesByTier

// LoadTemplates loads and parses templates from base dir
func LoadTemplates(source string) (Templates, error) {
	templates := make(map[ArtifactKind]TemplatesByTier)

	err := filepath.Walk(source, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walk dir%s: %v", path, err)
		}
		if f.IsDir() {
			return nil
		}
		
		filename := filepath.Base(path)
		if !strings.HasSuffix(filename, ".yaml") {
			return nil
		}

		templateName := strings.Split(filename, ".")[0]
		templateChunks := strings.Split(templateName,"-")

		artifactName := templateChunks[0]
		artifactTier := ""

		if len(templateChunks) > 1 {
			artifactTier = templateChunks[1]
		}

		var artifactKind ArtifactKind
		switch artifactName {
		case CronJobName: artifactKind = CronJobKind
		case DeploymentName: artifactKind = DeploymentKind
		case ServiceName: artifactKind = ServiceKind
		default: {
			print(color.HiBlackString("      unknown template type\n"))
			return nil
		}
		}

		filedata, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		parsedTmpl, err := template.New("Manifest").Parse(string(filedata[:len(filedata)]))
		if err != nil {
			return fmt.Errorf("can not parse template for cronjob manifest: %v", err)
		}
	
		tbt, ok := templates[artifactKind]
		if !ok {
			tbt = make(map[string]*template.Template)
			templates[artifactKind] = tbt
		}
		tbt[artifactTier] = parsedTmpl

		return nil
	})

	return templates, err
}
