package morpheus_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/spoonboy-io/dozer/internal/morpheus"
	"github.com/spoonboy-io/dozer/internal/state"
)

func TestGetProcesses(t *testing.T) {
	// test expectations
	wantSt := &state.State{
		LastPollProcessId:  4,
		ExecutingProcesses: []int{2, 4},
	}

	gotSt := &state.State{}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// add some mock rows
	rows := sqlmock.NewRows([]string{
		"id", "sub_type", "updated_by_id", "output_format", "date_created", "server_name",
		"created_by_id", "process_type_id", "updated_by", "updated_by_display_name", "error", "app_name", "success",
		"created_by_display_name", "display_name", "input", "app_id", "message", "ref_type", "job_template_id",
		"container_name", "output", "api_key", "account_id", "status_eta", "timer_sub_category", "process_type_name",
		"task_set_name", "container_id", "job_template_name", "task_set_id", "last_updated", "server_group_name", "sub_id",
		"deleted", "task_id", "unique_id", "percent", "timer_category", "reason", "end_date", "duration", "instance_name",
		"start_date", "zone_id", "input_format", "server_id", "exit_code", "integration_id", "ref_id", "instance_id",
		"server_group_id", "task_name", "created_by", "status", "process_result", "description", "event_title"}).
		AddRow(
			1, "", 1, "", time.Now(), "", 0, 0, "", "", "", "", false, "", "", "", 0, "", "", 0, "", "", "",
			0, 0, "", "", "", 0, "", 0, time.Now(), "", 0, "", 0, "", 0, "", "", time.Now(), 0, "", time.Now(), 0, "",
			0, "", 0, 0, 0, 0, "", "", "complete", "", "", "",
		).
		AddRow(
			2, "", 1, "", time.Now(), "", 0, 0, "", "", "", "", false, "", "", "", 0, "", "", 0, "", "", "",
			0, 0, "", "", "", 0, "", 0, time.Now(), "", 0, "", 0, "", 0, "", "", time.Now(), 0, "", time.Now(), 0, "",
			0, "", 0, 0, 0, 0, "", "", "executing", "", "", "",
		).
		AddRow(
			3, "", 1, "", time.Now(), "", 0, 0, "", "", "", "", false, "", "", "", 0, "", "", 0, "", "", "",
			0, 0, "", "", "", 0, "", 0, time.Now(), "", 0, "", 0, "", 0, "", "", time.Now(), 0, "", time.Now(), 0, "",
			0, "", 0, 0, 0, 0, "", "", "complete", "", "", "",
		).
		AddRow(
			4, "", 1, "", time.Now(), "", 0, 0, "", "", "", "", false, "", "", "", 0, "", "", 0, "", "", "",
			0, 0, "", "", "", 0, "", 0, time.Now(), "", 0, "", 0, "", 0, "", "", time.Now(), 0, "", time.Now(), 0, "",
			0, "", 0, 0, 0, 0, "", "", "executing", "", "", "",
		)

	mock.ExpectQuery("^SELECT (.+) FROM process where id > (.+)*").WillReturnRows(rows)

	if err := morpheus.GetProcesses(db, gotSt); err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// check state last process ID
	if gotSt.LastPollProcessId != wantSt.LastPollProcessId {
		t.Errorf("failed got %v wanted %v", gotSt.LastPollProcessId, wantSt.LastPollProcessId)
	}

	// check state executing processes
	if !reflect.DeepEqual(gotSt.ExecutingProcesses, wantSt.ExecutingProcesses) {
		t.Errorf("failed got %v wanted %v", gotSt.ExecutingProcesses, wantSt.ExecutingProcesses)
	}

	// check expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCheckExecuting(t *testing.T) {
	// test expectations
	wantSt := &state.State{
		LastPollProcessId:  4,
		ExecutingProcesses: []int{4},
	}

	gotSt := &state.State{
		ExecutingProcesses: []int{2, 4},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// add some mock rows
	// process id 2 is now complete so should be removed from state.ExecutingProcesses
	rows := sqlmock.NewRows([]string{
		"id", "sub_type", "updated_by_id", "output_format", "date_created", "server_name",
		"created_by_id", "process_type_id", "updated_by", "updated_by_display_name", "error", "app_name", "success",
		"created_by_display_name", "display_name", "input", "app_id", "message", "ref_type", "job_template_id",
		"container_name", "output", "api_key", "account_id", "status_eta", "timer_sub_category", "process_type_name",
		"task_set_name", "container_id", "job_template_name", "task_set_id", "last_updated", "server_group_name", "sub_id",
		"deleted", "task_id", "unique_id", "percent", "timer_category", "reason", "end_date", "duration", "instance_name",
		"start_date", "zone_id", "input_format", "server_id", "exit_code", "integration_id", "ref_id", "instance_id",
		"server_group_id", "task_name", "created_by", "status", "process_result", "description", "event_title"}).
		AddRow(
			2, "", 1, "", time.Now(), "", 0, 0, "", "", "", "", false, "", "", "", 0, "", "", 0, "", "", "",
			0, 0, "", "", "", 0, "", 0, time.Now(), "", 0, "", 0, "", 0, "", "", time.Now(), 0, "", time.Now(), 0, "",
			0, "", 0, 0, 0, 0, "", "", "complete", "", "", "",
		).
		AddRow(
			4, "", 1, "", time.Now(), "", 0, 0, "", "", "", "", false, "", "", "", 0, "", "", 0, "", "", "",
			0, 0, "", "", "", 0, "", 0, time.Now(), "", 0, "", 0, "", 0, "", "", time.Now(), 0, "", time.Now(), 0, "",
			0, "", 0, 0, 0, 0, "", "", "executing", "", "", "",
		)

	mock.ExpectQuery("^SELECT (.+) FROM process where id in \\((.+)\\)*").WillReturnRows(rows)

	if err := morpheus.CheckExecuting(db, gotSt); err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// check state executing processes
	if !reflect.DeepEqual(gotSt.ExecutingProcesses, wantSt.ExecutingProcesses) {
		t.Errorf("failed got %v wanted %v", gotSt.ExecutingProcesses, wantSt.ExecutingProcesses)
	}

	// check expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
