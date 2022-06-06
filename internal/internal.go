package internal

import (
	"database/sql"
)

const (
	POLL_INTERVAL = 5
)

// ProcessType is a struct to represent a morpheus process type
// image_name and description are omitted as they appear to be unused
type ProcessType struct {
	Id        int            `db:"id"`
	Code      sql.NullString `db:"code"`
	Name      sql.NullString `db:"name"`
	ImageCode sql.NullString `db:"image_code"`
}

// ProcessTypes uses process_type code as key and name as value
// so we can use code in the YAML config but look up against name in the process table
var ProcessTypes map[string]string

// Process is a struct to represent a morpheus process and all the possible
// information reported by morpheus about the process
type Process struct {
	Id                   int             `db:"id"`
	SubType              sql.NullString  `db:"sub_type"`
	UpdatedById          sql.NullInt64   `db:"updated_by_id"`
	OutputFormat         sql.NullString  `db:"output_format"`
	DateCreated          sql.NullTime    `db:"date_created"`
	ServerName           sql.NullString  `db:"server_name"`
	CreatedById          sql.NullInt64   `db:"created_by_id"`
	ProcessTypeId        sql.NullInt64   `db:"process_type_id"`
	UpdatedBy            sql.NullString  `db:"updated_by"`
	UpdatedByDisplayName sql.NullString  `db:"updated_by_display_name"`
	Error                sql.NullString  `db:"error"`
	AppName              sql.NullString  `db:"app_name"`
	Success              sql.NullBool    `db:"success"`
	CreatedByDisplayName sql.NullString  `db:"created_by_display_name"`
	DisplayName          sql.NullString  `db:"display_name"`
	Input                sql.NullString  `db:"input"`
	AppId                sql.NullInt64   `db:"app_id"`
	Message              sql.NullString  `db:"message"`
	RefType              sql.NullString  `db:"ref_type"`
	JobTemplateId        sql.NullInt64   `db:"job_template_id"`
	ContainerName        sql.NullString  `db:"container_name"`
	Output               sql.NullString  `db:"output"`
	ApiKey               sql.NullString  `db:"api_key"`
	AccountId            sql.NullInt64   `db:"account_id"`
	StatusEta            sql.NullInt64   `db:"status_eta"`
	TimerSubCategory     sql.NullString  `db:"timer_sub_category"`
	ProcessTypeName      sql.NullString  `db:"process_type_name"`
	TaskSetName          sql.NullString  `db:"task_set_name"`
	ContainerId          sql.NullInt64   `db:"container_id"`
	JobTemplateName      sql.NullString  `db:"job_template_name"`
	TaskSetId            sql.NullInt64   `db:"task_set_id"`
	LastUpdated          sql.NullTime    `db:"last_updated"`
	ServerGroupName      sql.NullString  `db:"server_group_name"`
	SubId                sql.NullInt64   `db:"sub_id"`
	Deleted              []byte          `db:"deleted"` // TODO temporary, could trip us up, it's a BOOL really
	TaskId               sql.NullInt64   `db:"task_id"`
	UniqueId             sql.NullString  `db:"unique_id"`
	Percent              sql.NullFloat64 `db:"percent"`
	TimerCategory        sql.NullString  `db:"timer_category"`
	Reason               sql.NullString  `db:"reason"`
	EndDate              sql.NullTime    `db:"end_date"`
	Duration             sql.NullInt64   `db:"duration"`
	InstanceName         sql.NullString  `db:"instance_name"`
	StartDate            sql.NullTime    `db:"start_date""`
	ZoneId               sql.NullInt64   `db:"zone_id"`
	InputFormat          sql.NullString  `db:"input_format"`
	ServerId             sql.NullInt64   `db:"server_id"`
	ExitCode             sql.NullString  `db:"exit_code"`
	IntegrationId        sql.NullInt64   `db:"integration_id"`
	RefId                sql.NullInt64   `db:"ref_id"`
	InstanceId           sql.NullInt64   `db:"instance_id"`
	ServergroupId        sql.NullInt64   `db:"server_group_id"`
	TaskName             sql.NullString  `db:"task_name"`
	CreatedBy            sql.NullString  `db:"created_by"`
	Status               string          `db:"status"`
	ProcessResult        sql.NullString  `db:"process_result"`
	Description          sql.NullString  `db:"description"`
	EventTitle           sql.NullString  `db:"event_title"`
}
