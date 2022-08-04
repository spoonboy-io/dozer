package hook

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"

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

// adding some standard errors we can check in the tests
// specifically for validating the config
var (
	ERR_NO_DESCRIPTION              = errors.New("No description is set")
	ERR_BAD_METHOD                  = errors.New("method is not acceptable")
	ERR_BAD_URL                     = errors.New("url is appears to be invalid")
	ERR_NO_BODY                     = errors.New("method requires requestBody")
	ERR_NO_TRIGGER                  = errors.New("No triggers defined in the hook")
	ERR_BAD_STATUS_TRIGGER          = errors.New("Trigger set on status is not recognised")
	ERR_NO_EXECUTING_STATUS_TRIGGER = errors.New("Can not trigger on status 'executing'")
	ERR_NOT_HTTPS                   = errors.New("url is not secure (no HTTPS)")
	ERR_COULD_NOT_PARSE_BODY        = errors.New("Problem parsing request body, check included variables")
)

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

func ValidateConfig() error {
	for i := range config {
		// check description
		if config[i].Description == "" {
			return ERR_NO_DESCRIPTION
		}

		// check method
		if err := isGoodMethod(config[i].Method); err != nil {
			return err
		}

		// check url
		if _, err := url.ParseRequestURI(config[i].URL); err != nil {
			return ERR_BAD_URL
		}

		// Reference: https://github.com/spoonboy-io/dozer/issues/1
		// temporary removal of validation
		/*
			else if pURL.Scheme != "https" {
				return ERR_NOT_HTTPS
			}
		*/

		// if method POST/PUT check request body is present
		if err := shouldHaveRequestBody(config[i].Method, config[i].RequestBody); err != nil {
			return err
		}

		// check at least one trigger
		if err := checkTriggers(config[i].Triggers); err != nil {
			return err
		}
	}

	return nil
}

// helpers
func isGoodMethod(method string) error {
	switch method {
	case "GET", "POST":
		return nil
	default:
		return ERR_BAD_METHOD
	}
}

func shouldHaveRequestBody(method, requestBody string) error {
	if method != "GET" {
		if requestBody == "" {
			return ERR_NO_BODY
		}
		// we should parse the body, to check that any included vars are valid
		// or we'll fail at runtime when the hook is trigger

		_, err := parseRequestBody(&internal.Process{}, requestBody)
		if err != nil {
			return ERR_COULD_NOT_PARSE_BODY
		}
	}
	return nil
}

func checkTriggers(trigger Trigger) error {
	if trigger.TaskName == "" && trigger.Status == "" && trigger.ProcessType == "" && trigger.CreatedBy == "" && trigger.AccountId == 0 {
		return ERR_NO_TRIGGER
	}

	if trigger.Status == "executing" {
		return ERR_NO_EXECUTING_STATUS_TRIGGER
	}

	switch trigger.Status {
	case "complete", "failed", "":
		return nil
	default:
		fmt.Println(trigger.Status)
		return ERR_BAD_STATUS_TRIGGER
	}
}
