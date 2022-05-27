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
	Hook struct {
		URL      string `yaml:"url"`
		Method   string `yaml:"method"`
		Token    string `yaml:"token"`
		Triggers struct {
			Status      string `yaml:"status"`
			ProcessType string `yaml:"processType"`
			TaskName    string `yaml:"taskName"`
			AccountId   int    `yaml:"accountId"`
			CreatedBy   string `yaml:"createdBy"`
		} `yaml:"triggers"`
	} `yaml:"webhook"`
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

// CheckProcess will check a process against the configuration to determine if
// it is an event that should trigger a call webhook
func CheckProcess(process internal.Process) error {
	fmt.Println("checking")

	return nil
}
