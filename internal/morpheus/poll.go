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
	EXECUTING = "executing"
)

// GetProcesses polls the database for processes higher than the store latestProcessId
// if the process is found to be executing it will be tracked, otherwise it is passed on for
// checking against the webhook configuration
func GetProcesses(db *sql.DB, st *state.State, logger *koan.Logger, ctx context.Context) error {
	rows, err := db.Query("SELECT * FROM process where id > ?;", st.LastPollProcessId)
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
			go hook.CheckProcess(&process, logger, ctx)
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
func CheckExecuting(db *sql.DB, st *state.State, logger *koan.Logger, ctx context.Context) error {
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
	query := fmt.Sprintf("SELECT * FROM process where id in (%s);", processList)
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
			go hook.CheckProcess(&process, logger, ctx)
			// delete from state
			st.DeleteProcessFromState(process.Id)
		}
	}

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
		&process.ServerID,
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
