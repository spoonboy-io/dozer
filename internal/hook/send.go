package hook

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/spoonboy-io/koan"

	"github.com/spoonboy-io/dozer/internal"
)

// safeProcess will collect the properties we want to make available for interpolating into the JSON request body,
// we make available via this type and not process as the sql null int64, string and bool types are not parseable
type safeProcess struct {
	Id                   int
	SubType              string
	UpdatedById          int64
	OutputFormat         string
	DateCreated          time.Time
	ServerName           string
	CreatedById          int64
	ProcessTypeId        int64
	UpdatedBy            int64
	UpdatedByDisplayName string
	Error                string
	AppName              string
	Success              bool
	CreatedByDisplayName string
	DisplayName          string
	Input                string
	AppId                int64
	Message              string
	RefType              string
	JobTemplateId        int64
	ContainerName        string
	Output               string
	ApiKey               string
	AccountId            int64
	StatusEta            int64
	TimerSubCategory     string
	ProcessTypeName      string
	TaskSetName          string
	ContainerId          int64
	JobTemplateName      string
	TaskSetId            int64
	LastUpdated          time.Time
	ServerGroupName      string
	SubId                int64
	Deleted              []byte
	TaskId               int64
	UniqueId             string
	Percent              float64
	TimerCategory        string
	Reason               string
	EndDate              time.Time
	Duration             int64
	InstanceName         string
	StartDate            time.Time
	ZoneId               string
	InputFormat          string
	ServerID             string
	ExitCode             string
	IntegrationId        int64
	RefId                int64
	InstanceId           int64
	ServergroupId        int64
	TaskName             string
	CreatedBy            string
	Status               string
	ProcessResult        string
	Description          string
	EventTitle           string
}

func fireWebhook(ctx context.Context, process *internal.Process, hook *Hook, logger *koan.Logger) error {
	var data io.Reader
	var err error

	// parse RequestBody if required
	if hook.Method != "GET" && hook.RequestBody != "" {
		data, err = parseRequestBody(process, hook.RequestBody)
		if err != nil {
			return err
		}
	}

	// form the request, make and return any errors
	req, err := http.NewRequest(hook.Method, hook.URL, data)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	// form the authorization header if exists
	if hook.Token != "" {
		req.Header.Add("Authorization", hook.Token)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Bad response (%d): Hook: %s, URL: %s", res.StatusCode, hook.Description, hook.URL)
	}

	return nil
}

func parseRequestBody(process *internal.Process, body string) (io.Reader, error) {
	var buffer bytes.Buffer

	safeProcess := safeProcess{
		Id:        process.Id,
		SubType:   process.SubType.String,
		UpdatedBy: process.UpdatedById.Int64,
		// TODO we need to add in the remaining properties
	}

	fmt.Println("before: ", body)

	t := template.Must(template.New("body").Parse(body))
	if err := t.Execute(&buffer, safeProcess); err != nil {
		return &buffer, err
	}

	fmt.Println("after: ", buffer.String())

	return &buffer, nil
}
