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
	Description string `yaml:"description"`
	URL         string `yaml:"url"`
	Method      string `yaml:"method"`
	Token       string `yaml:"token"`
	Triggers    struct {
		Status      string `yaml:"status"`
		ProcessType string `yaml:"processType"`
		TaskName    string `yaml:"taskName"`
		AccountId   int    `yaml:"accountId"`
		CreatedBy   string `yaml:"createdBy"`
	} `yaml:"triggers"`
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
func CheckProcess(process internal.Process) {
	// go through all the hook config
	for _, h := range config {
		// config uses code we need to search on name, we swap the code for the name in the config
		processName, err := getProcessTypeName(h.Triggers.ProcessType)
		if err != nil {
			// TODO we choose to move on for now
			continue
		}
		h.Triggers.ProcessType = processName

		if fire := checkStatus(process, h.Hook); fire {
			fmt.Println("firing on status", process.Id, process.Status)
			continue
		}

		if fire := checkProcessType(process, h.Hook); fire {
			fmt.Println("firing on process type", process.Id, process.ProcessTypeName.String)
			continue
		}

		if fire := checkTaskName(process, h.Hook); fire {
			fmt.Println("firing on task name", process.Id, process.TaskName.String)
			continue
		}

		if fire := checkAccountId(process, h.Hook); fire {
			fmt.Println("firing on account id", process.Id, process.AccountId.Int64)
			continue
		}

		if fire := checkCreatedBy(process, h.Hook); fire {
			fmt.Println("firing on created by", process.Id, process.CreatedBy.String)
			continue
		}
	}
}

func checkStatus(process internal.Process, hook Hook) bool {
	var fire bool
	if hook.Triggers.Status == process.Status {
		fire = true
		if hook.Triggers.ProcessType != "" {
			fire = false
			if hook.Triggers.ProcessType == process.ProcessTypeName.String {
				fire = true
			}
		}
		if hook.Triggers.TaskName != "" {
			fire = false
			if hook.Triggers.TaskName == process.TaskName.String {
				fire = true
			}
		}

		if hook.Triggers.AccountId != 0 {
			fire = false
			if hook.Triggers.AccountId == int(process.AccountId.Int64) {
				fire = true
			}
		}

		if hook.Triggers.CreatedBy != "" {
			fire = false
			if hook.Triggers.CreatedBy == process.CreatedBy.String {
				fire = true
			}
		}
	}
	return fire
}

func checkProcessType(process internal.Process, hook Hook) bool {
	var fire bool
	if hook.Triggers.ProcessType == process.ProcessTypeName.String {
		fire = true
		if hook.Triggers.Status != "" {
			fire = false
			if hook.Triggers.Status == process.Status {
				fire = true
			}
		}
		if hook.Triggers.TaskName != "" {
			fire = false
			if hook.Triggers.TaskName == process.TaskName.String {
				fire = true
			}
		}

		if hook.Triggers.AccountId != 0 {
			fire = false
			if hook.Triggers.AccountId == int(process.AccountId.Int64) {
				fire = true
			}
		}

		if hook.Triggers.CreatedBy != "" {
			fire = false
			if hook.Triggers.CreatedBy == process.CreatedBy.String {
				fire = true
			}
		}
	}
	return fire
}

func checkTaskName(process internal.Process, hook Hook) bool {
	var fire bool
	if hook.Triggers.TaskName == process.TaskName.String {
		fire = true
		if hook.Triggers.Status != "" {
			fire = false
			if hook.Triggers.Status == process.Status {
				fire = true
			}
		}
		if hook.Triggers.ProcessType != "" {
			fire = false
			if hook.Triggers.ProcessType == process.ProcessTypeName.String {
				fire = true
			}
		}

		if hook.Triggers.AccountId != 0 {
			fire = false
			if hook.Triggers.AccountId == int(process.AccountId.Int64) {
				fire = true
			}
		}

		if hook.Triggers.CreatedBy != "" {
			fire = false
			if hook.Triggers.CreatedBy == process.CreatedBy.String {
				fire = true
			}
		}
	}
	return fire
}

func checkAccountId(process internal.Process, hook Hook) bool {
	var fire bool
	if hook.Triggers.AccountId == int(process.AccountId.Int64) {
		fire = true
		if hook.Triggers.ProcessType != "" {
			fire = false
			if hook.Triggers.ProcessType == process.ProcessTypeName.String {
				fire = true
			}
		}
		if hook.Triggers.TaskName != "" {
			fire = false
			if hook.Triggers.TaskName == process.TaskName.String {
				fire = true
			}
		}

		if hook.Triggers.Status != "" {
			fire = false
			if hook.Triggers.Status == process.Status {
				fire = true
			}
		}

		if hook.Triggers.CreatedBy != "" {
			fire = false
			if hook.Triggers.CreatedBy == process.CreatedBy.String {
				fire = true
			}
		}
	}
	return fire
}

func checkCreatedBy(process internal.Process, hook Hook) bool {
	var fire bool
	if hook.Triggers.CreatedBy == process.CreatedBy.String {
		fire = true
		if hook.Triggers.ProcessType != "" {
			fire = false
			if hook.Triggers.ProcessType == process.ProcessTypeName.String {
				fire = true
			}
		}
		if hook.Triggers.TaskName != "" {
			fire = false
			if hook.Triggers.TaskName == process.TaskName.String {
				fire = true
			}
		}

		if hook.Triggers.AccountId != 0 {
			fire = false
			if hook.Triggers.AccountId == int(process.AccountId.Int64) {
				fire = true
			}
		}

		if hook.Triggers.Status != "" {
			fire = false
			if hook.Triggers.Status == process.Status {
				fire = true
			}
		}
	}

	return fire
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
