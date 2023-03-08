package morpheus

import "testing"

func TestEasyQuery(t *testing.T) {
	testCases := []struct {
		name   string
		suffix string
		want   string
	}{
		{
			"use 1",
			"> ?",
			`SELECT id, sub_type, updated_by_id, output_format, date_created, server_name, created_by_id, process_type_id, updated_by, updated_by_display_name,
error, app_name, success, created_by_display_name, display_name, input, app_id, message, ref_type, job_template_id, container_name, output,
api_key, account_id, status_eta, timer_sub_category, process_type_name, task_set_name, container_id, job_template_name, task_set_id, last_updated, server_group_name,
sub_id, deleted, task_id, unique_id, percent, timer_category, reason, end_date, duration, instance_name, start_date, zone_id, input_format, server_id,
exit_code, integration_id, ref_id, instance_id, server_group_id, task_name, created_by, status, process_result, description, event_title FROM process where id > ?;`,
		},
		{
			"use 2",
			"in (%s)",
			`SELECT id, sub_type, updated_by_id, output_format, date_created, server_name, created_by_id, process_type_id, updated_by, updated_by_display_name,
error, app_name, success, created_by_display_name, display_name, input, app_id, message, ref_type, job_template_id, container_name, output,
api_key, account_id, status_eta, timer_sub_category, process_type_name, task_set_name, container_id, job_template_name, task_set_id, last_updated, server_group_name,
sub_id, deleted, task_id, unique_id, percent, timer_category, reason, end_date, duration, instance_name, start_date, zone_id, input_format, server_id,
exit_code, integration_id, ref_id, instance_id, server_group_id, task_name, created_by, status, process_result, description, event_title FROM process where id in (%s);`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := easyQuery(tc.suffix)
			if got != tc.want {
				t.Errorf("wanted %v got %v", tc.want, got)
			}
		})
	}

}
