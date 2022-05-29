package morpheus

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spoonboy-io/dozer/internal/hook"

	"github.com/spoonboy-io/dozer/internal"
	"github.com/spoonboy-io/dozer/internal/state"
)

const (
	COMPLETE  = "complete"
	FAILED    = "failed"
	EXECUTING = "executing"
)

func GetProcesses(db *sql.DB, st *state.State) error {
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
			// TODO potentially we should use goroutine so we don't block
			hook.CheckProcess(process)
		}
		lastProcessId = process.Id
		//fmt.Printf("Process: %+v\n\n\n", process)
	}

	// update state
	st.LastPollTimestamp = time.Now().Round(0)
	if lastProcessId > st.LastPollProcessId {
		st.LastPollProcessId = lastProcessId
	}

	return nil
}

func CheckExecuting(db *sql.DB, st *state.State) error {
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
		fmt.Println(process)
		if process.Status != EXECUTING {
			// status is complete or failed so compare row to hook configuration
			hook.CheckProcess(process)
			// delete from state
			st.DeleteProcessFromState(process.Id)
		}
	}

	return nil
}

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
