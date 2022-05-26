package internal

import "time"

const (
	POLL_INTERVAL = 5
)

// ProcessType is a struct to represent a morpheus process type
// image_name and description are omitted as they appear to be unused
type ProcessType struct {
	Id        int    `db:"id"`
	Code      string `db:"code"`
	Name      string `db:"name"`
	ImageCode string `db:"image_code"`
}

// Process is a struct to represent a morpheus process and all the possible
// information reported by morpheus about the process
type Process struct {
	Id                   int       `db:"id"`
	SubType              string    `db:"sub_type"`
	UpdatedById          int       `db:"updated_by_id"`
	OutputFormat         string    `db:"output_format"`
	DateCreated          time.Time `db:"date_created"`
	ServerName           string    `db:"server_name"`
	CreatedById          int       `db:"created_by_id"`
	ProcessTypeId        int       `db:"process_type_id"`
	UpdatedBy            string    `db:"updated_by"`
	UpdatedByDisplayName string    `db:"updated_by_display_name"`
	Error                string    `db:"error"`
	AppName              string    `db:"app_name"`
	Succeess             bool      `db:"success""`
	CreateByDisplayName  string    `db:"created_by_display_name"`
	DisplayName          string    `db:"display_name"`
	Input                string    `db:"input"`
	AppId                int       `db:"app_id"`
	Message              string    `db:"message"`
	RefType              string    `db:"ref_type"`
	JobTemplateId        int       `db:"job_template_id"`
	ContainerName        string    `db:"container_name"`
	Output               string    `db:"output"`
	ApiKey               string    `db:"api_key"`
	AccountId            int       `db:"account_id"`
	StatusEta            int       `db:"status_eta"`
	TimerSubCategory     string    `db:"timer_sub_category"`
	ProcessTypeName      string    `db:"process_type_name"`
	TaskSetName          string    `db:"task_set_name"`
	ContainerId          int       `db:"container_id"`
	JobTemplateName      string    `db:"job_template_name"`
	TaskSetId            int       `db:"task_set_id"`
	LastUpdated          time.Time `db:"last_updated"`
	ServerGroupName      string    `db:"server_group_name"`
	SubId                int       `db:"sub_id"`
	Deleted              bool      `db:"deleted"`
	TaskId               int       `db:"task_id"`
	UniqueId             string    `db:"unique_id"`
	Percent              float32   `db:"percent"`
	TimerCategory        string    `db:"timer_category"`
	Reason               string    `db:"reason"`
	EndDate              time.Time `db:"end_date"`
	Duration             int       `db:"duration"`
	InstanceName         string    `db:"instance_name"`
	StartDate            time.Time `db:"start_date""`
	ZoneId               int       `db:"zone_id"`
	InputFormat          string    `db:"input_format"`
	ServerID             int       `db:"server_id"`
	ExitCode             string    `db:"exit_code"`
	IntegrationId        int       `db:"integration_id"`
	RefId                int       `db:"ref_id"`
	InstanceId           int       `db:"instance_id"`
	ServergroupId        int       `db:"server_group_id"`
	TaskName             string    `db:"task_name"`
	CreatedBy            string    `db:"created_by"`
	Status               string    `db:"status"`
	ProcessResult        string    `db:"process_result"`
	Description          string    `db:"description"`
	EventTitle           string    `db:"event_title"`
}
