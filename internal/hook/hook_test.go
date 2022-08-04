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

func TestValidateConfig(t *testing.T) {
	testCases := []struct {
		name    string
		config  Hooks
		wantErr error
	}{
		{
			name: "all good, should pass",
			config: Hooks{
				{
					Hook{
						Description: "test hook 1",
						URL:         "https://testurl.com",
						Method:      "GET",
						Token:       "faketoken1",
						Triggers: Trigger{
							Status: "complete",
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "bad method, should fail",
			config: Hooks{
				{
					Hook{
						Description: "test hook 1",
						URL:         "https://testurl.com",
						Method:      "BAD_GET",
						Token:       "faketoken1",
						Triggers: Trigger{
							Status: "complete",
						},
					},
				},
			},
			wantErr: ERR_BAD_METHOD,
		},
		{
			name: "no description, should fail",
			config: Hooks{
				{
					Hook{
						Description: "",
						URL:         "https://testurl.com",
						Method:      "GET",
						Token:       "faketoken1",
						Triggers: Trigger{
							Status: "complete",
						},
					},
				},
			},
			wantErr: ERR_NO_DESCRIPTION,
		},
		{
			name: "bad url, should fail",
			config: Hooks{
				{
					Hook{
						Description: "test hook",
						URL:         "https//testurl.com",
						Method:      "GET",
						Token:       "faketoken1",
						Triggers: Trigger{
							Status: "complete",
						},
					},
				},
			},
			wantErr: ERR_BAD_URL,
		},
		{
			name: "no request body on POST, should fail",
			config: Hooks{
				{
					Hook{
						Description: "test hook",
						URL:         "https://testurl.com",
						Method:      "POST",
						Token:       "faketoken1",
						Triggers: Trigger{
							Status: "complete",
						},
					},
				},
			},
			wantErr: ERR_NO_BODY,
		},
		{
			name: "bad request body variable {{.BadId}}, should fail",
			config: Hooks{
				{
					Hook{
						Description: "test hook",
						URL:         "https://testurl.com",
						Method:      "POST",
						RequestBody: "{{.BadId}}",
						Token:       "faketoken1",
						Triggers: Trigger{
							Status: "complete",
						},
					},
				},
			},
			wantErr: ERR_COULD_NOT_PARSE_BODY,
		},
		{
			name: "no triggers are included, should fail",
			config: Hooks{
				{
					Hook{
						Description: "test hook",
						URL:         "https://testurl.com",
						Method:      "GET",
						Token:       "faketoken1",
					},
				},
			},
			wantErr: ERR_NO_TRIGGER,
		},
		{
			name: "status trigger is `executing` should fail",
			config: Hooks{
				{
					Hook{
						Description: "test hook",
						URL:         "https://testurl.com",
						Method:      "GET",
						Token:       "faketoken1",
						Triggers: Trigger{
							Status: "executing",
						},
					},
				},
			},
			wantErr: ERR_NO_EXECUTING_STATUS_TRIGGER,
		},
		{
			name: "status trigger invalid/mispelt should fail",
			config: Hooks{
				{
					Hook{
						Description: "test hook",
						URL:         "https://testurl.com",
						Method:      "GET",
						Token:       "faketoken1",
						Triggers: Trigger{
							Status: "badcomplete",
						},
					},
				},
			},
			wantErr: ERR_BAD_STATUS_TRIGGER,
		},

		// Reference: https://github.com/spoonboy-io/dozer/issues/1
		// temporary removal of validation
		/*
			{
				name: "url is not HTTPS, should fail",
				config: Hooks{
					{
						Hook{
							Description: "test hook 1",
							URL:         "http://testurl.com",
							Method:      "GET",
							Triggers: Trigger{
								Status: "complete",
							},
						},
					},
				},
				wantErr: ERR_NOT_HTTPS,
			},
		*/
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// set the package config
			config = tc.config
			gotErr := ValidateConfig()
			if gotErr != tc.wantErr {
				t.Errorf("wanted %v got %v", tc.wantErr, gotErr)
			}
		})
	}
}
