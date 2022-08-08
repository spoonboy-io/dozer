package hook

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/spoonboy-io/dozer/internal"
)

// safeProcess will collect the properties we want to make available for interpolating into the JSON request body,
// we make available via this type and not process as the sql null int64, string and bool types are not parseable
// some properties are not available as they may leak secrets, or be inapplicable to complete or failed processes
type safeProcess struct {
	Id                   int
	SubType              string
	UpdatedById          int64
	OutputFormat         string
	DateCreated          time.Time
	ServerName           string
	CreatedById          int64
	ProcessTypeId        int64
	UpdatedBy            string
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
	//ApiKey               string
	AccountId int64
	//StatusEta int64
	//TimerSubCategory     string
	ProcessTypeName string
	TaskSetName     string
	ContainerId     int64
	JobTemplateName string
	TaskSetId       int64
	LastUpdated     time.Time
	ServerGroupName string
	SubId           int64
	//Deleted              []byte
	TaskId   int64
	UniqueId string
	//Percent              float64
	//TimerCategory        string
	Reason        string
	EndDate       time.Time
	Duration      int64
	InstanceName  string
	StartDate     time.Time
	ZoneId        int64
	InputFormat   string
	ServerId      int64
	ExitCode      string
	IntegrationId int64
	RefId         int64
	InstanceId    int64
	ServergroupId int64
	TaskName      string
	CreatedBy     string
	Status        string
	ProcessResult string
	Description   string
	EventTitle    string
}

func fireWebhook(ctx context.Context, process *internal.Process, hook *Hook) error {
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
	req.Header.Add("Content-Type", "application/json")

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
		Id:                   process.Id,
		SubType:              process.SubType.String,
		UpdatedById:          process.UpdatedById.Int64,
		OutputFormat:         process.OutputFormat.String,
		DateCreated:          process.DateCreated.Time,
		ServerName:           process.ServerName.String,
		CreatedById:          process.CreatedById.Int64,
		ProcessTypeId:        process.ProcessTypeId.Int64,
		UpdatedBy:            process.UpdatedBy.String,
		UpdatedByDisplayName: process.UpdatedByDisplayName.String,
		Error:                process.Error.String,
		AppName:              process.AppName.String,
		Success:              process.Success.Bool,
		CreatedByDisplayName: process.CreatedByDisplayName.String,
		DisplayName:          process.DisplayName.String,
		Input:                process.Input.String,
		AppId:                process.AppId.Int64,
		Message:              process.Message.String,
		RefType:              process.RefType.String,
		JobTemplateId:        process.JobTemplateId.Int64,
		ContainerName:        process.ContainerName.String,
		Output:               process.Output.String,
		AccountId:            process.AccountId.Int64,
		ProcessTypeName:      process.ProcessTypeName.String,
		TaskSetName:          process.TaskSetName.String,
		ContainerId:          process.ContainerId.Int64,
		JobTemplateName:      process.JobTemplateName.String,
		TaskSetId:            process.TaskSetId.Int64,
		LastUpdated:          process.LastUpdated.Time,
		ServerGroupName:      process.ServerGroupName.String,
		SubId:                process.SubId.Int64,
		TaskId:               process.TaskId.Int64,
		UniqueId:             process.UniqueId.String,
		Reason:               process.Reason.String,
		EndDate:              process.EndDate.Time,
		Duration:             process.Duration.Int64,
		InstanceName:         process.InstanceName.String,
		StartDate:            process.StartDate.Time,
		ZoneId:               process.ZoneId.Int64,
		InputFormat:          process.InputFormat.String,
		ServerId:             process.ServerId.Int64,
		ExitCode:             process.ExitCode.String,
		IntegrationId:        process.IntegrationId.Int64,
		RefId:                process.RefId.Int64,
		InstanceId:           process.InstanceId.Int64,
		ServergroupId:        process.ServergroupId.Int64,
		TaskName:             process.TaskName.String,
		CreatedBy:            process.CreatedBy.String,
		Status:               process.Status,
		ProcessResult:        process.ProcessResult.String,
		Description:          process.Description.String,
		EventTitle:           process.EventTitle.String,
	}

	t := template.Must(template.New("body").Parse(body))
	if err := t.Execute(&buffer, safeProcess); err != nil {
		return &buffer, err
	}

	return &buffer, nil
}
