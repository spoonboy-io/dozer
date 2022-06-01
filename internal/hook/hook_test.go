package hook

import (
	"os"
	"reflect"
	"testing"
)

var testYamlFile = "testWebhook.yaml"

func writeTestYamlFile(t *testing.T) {
	data := `---
- webhook:
    description: test hook 1
    url: https://testurl1
    method: GET
    token: faketoken1
    triggers:
      status: complete

- webhook:
    description: test hook 2
    url: https://testurl2
    method: POST
    requestBody: '{"id": 55}'
    token: faketoken2
    triggers:
      processType: provision`
	if err := os.WriteFile(testYamlFile, []byte(data), 0644); err != nil {
		t.Fatalf("could not write test yaml file %+v", err)
	}
}

func removeTestYamlFile(t *testing.T) {
	if err := os.Remove(testYamlFile); err != nil {
		t.Fatal("Could not remove tess yaml file")
	}
}

func TestReadAndParseConfig(t *testing.T) {
	// write test yaml config
	writeTestYamlFile(t)

	wantConfig := Hooks{
		{
			Hook{
				Description: "test hook 1",
				URL:         "https://testurl1",
				Method:      "GET",
				Token:       "faketoken1",
				Triggers: Trigger{
					Status: "complete",
				},
			},
		},
		{
			Hook{
				Description: "test hook 2",
				URL:         "https://testurl2",
				Method:      "POST",
				RequestBody: "{\"id\": 55}",
				Token:       "faketoken2",
				Triggers: Trigger{
					ProcessType: "provision",
				},
			},
		},
	}

	if err := ReadAndParseConfig(testYamlFile); err != nil {
		t.Fatalf("could not read test yaml file %+v", err)
	}

	gotConfig := config

	if !reflect.DeepEqual(gotConfig, wantConfig) {
		t.Errorf("\n\nWanted\n%v\n, \n\ngot \n%v\n", wantConfig, gotConfig)
	}

	// clean up
	removeTestYamlFile(t)
}
