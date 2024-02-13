package morpheus

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spoonboy-io/koan"

	"github.com/spoonboy-io/dozer/internal/hook"

	"github.com/spoonboy-io/dozer/internal"
	"github.com/spoonboy-io/dozer/internal/state"
)

const (
	EXECUTING = "running"
)

// GetProcesses polls the database for processes higher than the store latestProcessId
// if the process is found to be executing it will be tracked, otherwise it is passed on for
// checking against the webhook configuration
func GetProcesses(ctx context.Context, db *sql.DB, st *state.State, logger *koan.Logger) error {
	//rows, err := db.Query("SELECT * FROM process where id > ?;", st.LastPollProcessId)
	rows, err := db.Query(easyQuery("> ?"), st.LastPollProcessId)
	if err != nil {
		return err
	}

	var lastProcessId int
	for rows.Next() {
		var process internal.Process
		err := easyScan(&process, rows)
		if err != nil {
			return err
		}

		// track executing processes
		if process.Status == EXECUTING {
			st.ExecutingProcesses = append(st.ExecutingProcesses, process.Id)
		} else {
			// status is complete or failed so compare row to hook configuration
			go hook.CheckProcess(ctx, &process, logger)
		}
		lastProcessId = process.Id
	}

	// update state
	st.LastPollTimestamp = time.Now().Round(0)
	if lastProcessId > st.LastPollProcessId {
		st.LastPollProcessId = lastProcessId
	}

	return nil
}

// CheckExecuting uses state to obtain processes being tracked as 'executing', it performs
// an SQL query which checks their status. If found to be no longer in the executing state the
// process is passed on for checking against the webhook configuration and is no longer tracked in state
func CheckExecuting(ctx context.Context, db *sql.DB, st *state.State, logger *koan.Logger) error {
	if len(st.ExecutingProcesses) == 0 {
		// nothing to do
		return nil
	}

	// we build an IN sql query so it's one DB call
	var strList []string
	for _, v := range st.ExecutingProcesses {
		strList = append(strList, strconv.Itoa(v))
	}
	processList := strings.Join(strList, ",")
	//query := fmt.Sprintf("SELECT * FROM process where id in (%s);", processList)
	query := fmt.Sprintf(easyQuery("in (%s)"), processList)
	rows, err := db.Query(query)
	if err != nil {
		return err
	}

	for rows.Next() {
		var process internal.Process
		err := easyScan(&process, rows)
		if err != nil {
			return err
		}
		if process.Status != EXECUTING {
			// status is complete or failed so compare row to hook configuration
			go hook.CheckProcess(ctx, &process, logger)
			// delete from state
			st.DeleteProcessFromState(process.Id)
		}
	}

	return nil
}

// GetLastProcessIdOnStart is used to ascertain the latest process id on first run
// or where no previous state is available
func GetLastProcessIdOnStart(db *sql.DB, st *state.State) error {
	rows, err := db.Query("SELECT MAX(id) as max_id FROM process;")
	if err != nil {
		return err
	}

	var lastProcessId int
	for rows.Next() {
		err := rows.Scan(&lastProcessId)
		if err != nil {
			return err
		}
	}

	// update state
	st.LastPollProcessId = lastProcessId

	return nil
}

// just a helper to reduce the duplication
func easyScan(process *internal.Process, rows *sql.Rows) error {
	err := rows.Scan(
		&process.Id,
		&process.SubType,
		&process.UpdatedById,
		&process.OutputFormat,
		&process.DateCreated,
		&process.ServerName,
		&process.CreatedById,
		&process.ProcessTypeId,
		&process.UpdatedBy,
		&process.UpdatedByDisplayName,
		&process.Error,
		&process.AppName,
		&process.Success,
		&process.CreatedByDisplayName,
		&process.DisplayName,
		&process.Input,
		&process.AppId,
		&process.Message,
		&process.RefType,
		&process.JobTemplateId,
		&process.ContainerName,
		&process.Output,
		&process.ApiKey,
		&process.AccountId,
		&process.StatusEta,
		&process.TimerSubCategory,
		&process.ProcessTypeName,
		&process.TaskSetName,
		&process.ContainerId,
		&process.JobTemplateName,
		&process.TaskSetId,
		&process.LastUpdated,
		&process.ServerGroupName,
		&process.SubId,
		&process.Deleted,
		&process.TaskId,
		&process.UniqueId,
		&process.Percent,
		&process.TimerCategory,
		&process.Reason,
		&process.EndDate,
		&process.Duration,
		&process.InstanceName,
		&process.StartDate,
		&process.ZoneId,
		&process.InputFormat,
		&process.ServerId,
		&process.ExitCode,
		&process.IntegrationId,
		&process.RefId,
		&process.InstanceId,
		&process.ServergroupId,
		&process.TaskName,
		&process.CreatedBy,
		&process.Status,
		&process.ProcessResult,
		&process.Description,
		&process.EventTitle,
	)
	if err != nil {
		return err
	}
	return nil
}

// helper so we don't need to build this twice
func easyQuery(suffix string) string {
	q := `SELECT id, sub_type, updated_by_id, output_format, date_created, server_name, created_by_id, process_type_id, updated_by, updated_by_display_name,
error, app_name, success, created_by_display_name, display_name, input, app_id, message, ref_type, job_template_id, container_name, output,
api_key, account_id, status_eta, timer_sub_category, process_type_name, task_set_name, container_id, job_template_name, task_set_id, last_updated, server_group_name,
sub_id, deleted, task_id, unique_id, percent, timer_category, reason, end_date, duration, instance_name, start_date, zone_id, input_format, server_id,
exit_code, integration_id, ref_id, instance_id, server_group_id, task_name, created_by, status, process_result, description, event_title FROM process where id %s;`

	return fmt.Sprintf(q, suffix)
}
