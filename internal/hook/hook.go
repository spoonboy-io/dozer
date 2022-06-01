package hook

import (
	"fmt"
	"io/ioutil"

	"github.com/spoonboy-io/dozer/internal"

	"gopkg.in/yaml.v2"
)

var config Hooks

// Hooks is representation of the parsed YAML webhook configuration file
type Hooks []struct {
	Hook `yaml:"webhook"`
}

// Hook represents the configuration of a single webhook
type Hook struct {
	Description string  `yaml:"description"`
	URL         string  `yaml:"url"`
	Method      string  `yaml:"method"`
	Token       string  `yaml:"token"`
	RequestBody string  `yaml:"requestBody"`
	Triggers    Trigger `yaml:"triggers"`
}

// Trigger represents the trigger configuration options which can be set in the YAML.
// They are additive, in that all set, must be satisfied for the hook event to be fired
type Trigger struct {
	Status      string `yaml:"status"`
	ProcessType string `yaml:"processType"`
	TaskName    string `yaml:"taskName"`
	AccountId   int    `yaml:"accountId"`
	CreatedBy   string `yaml:"createdBy"`
}

// ReadAndParseConfig reads the contents of the YAML hook config filer
// and parses it to a map of Hook structs
func ReadAndParseConfig(cfgFile string) error {
	yamlConfig, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(yamlConfig, &config); err != nil {
		return err
	}
	return nil
}

// getProcessTypeName checks the key exists and returns it if found, or errors
func getProcessTypeName(code string) (string, error) {
	if code != "" {
		if _, ok := internal.ProcessTypes[code]; !ok {
			return "", fmt.Errorf("ProcessType not found, check YAML")
		}
		return internal.ProcessTypes[code], nil
	}
	return "", nil
}
