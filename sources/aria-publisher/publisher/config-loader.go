package publisher

import (
	"fmt"
	"os"
	"reflect"
)

// ServiceConfig config for image-builder
type ServiceConfig struct {
	CertFile       string `env:"CERT_FILE" default:"/certs/acc.io/acc.io.crt"`
	KeyFile        string `env:"KEY_FILE" default:"/certs/acc.io/acc.io.key"`
	ServerPort     string `env:"PORT" default:"9999"`
}

/* Non-exported instance to avoid accidental overwrite */
var serviceConfig ServiceConfig

/* Tag names to load configuration from environment variable */
const (
	ENV     = "env"
	DEFAULT = "default"
)

func loadFromEnv(tag string, defaultTag string) (string, string) {
	/* Check if the tag is defined in the environment or else replace with default value */
	envVar := os.Getenv(tag)
	if envVar == "" {
		envVar = defaultTag
		/* '1' is used to indicate that default value is being loaded */
		return envVar, "1"
	}
	return envVar, ""
}

// LoadConfiguration load configuration into config from environment variables
func LoadConfiguration() ServiceConfig {
	// ValueOf returns a Value representing the run-time data
	v := reflect.ValueOf(serviceConfig)
	for i := 0; i < v.NumField(); i++ {
		// Get the field tag value
		tag := v.Type().Field(i).Tag.Get(ENV)
		defaultTag := v.Type().Field(i).Tag.Get(DEFAULT)

		// Skip if tag is not defined or ignored
		if tag == "" || tag == "-" {
			continue
		}
		a := reflect.Indirect(reflect.ValueOf(serviceConfig))
		EnvVar, Info := loadFromEnv(tag, defaultTag)
		if Info != "" {
			fmt.Println("Missing environment configuration for '" + a.Type().Field(i).Name + "', Loading default setting!")
		}
		/* Set the value in the environment variable to the respective struct field */
		reflect.ValueOf(&serviceConfig).Elem().Field(i).SetString(EnvVar)
	}
	return serviceConfig
}
