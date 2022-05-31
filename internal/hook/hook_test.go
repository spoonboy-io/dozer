package hook

import (
	"database/sql"
	"testing"

	"github.com/spoonboy-io/dozer/internal"
)

// TestCheckProcessLogic contains testcases that checks that hooks fire correctly based on
// the configuration when test processes are inspected. We are not testing the webhook call itself
// we are whitebox testing as we are calling unexported functions
func TestCheckProcessesLogic(t *testing.T) {
	testCases := []struct {
		name                string
		process             internal.Process
		hook                Hook
		wantFireStatus      bool
		wantFireProcessType bool
		wantFireTaskName    bool
		wantFireAccountId   bool
		wantFireCreatedBy   bool
	}{
		// status alone
		{
			name: "Will fire on status complete",
			process: internal.Process{
				Status: "complete",
			},
			hook: Hook{
				Triggers: Trigger{
					Status: "complete",
				},
			},
			wantFireStatus: true,
		},
		{
			name: "Will not fire on status complete",
			process: internal.Process{
				Status: "complete",
			},
			hook: Hook{
				Triggers: Trigger{
					Status: "failed",
				},
			},
			wantFireStatus: false,
		},
		{
			name: "Will not fire on status executing",
			process: internal.Process{
				Status: "executing",
			},
			hook: Hook{
				Triggers: Trigger{
					Status: "executing",
				},
			},
			wantFireStatus: false,
		},

		// processType alone
		{
			name: "Fires for processTypeName 'local workflow' ",
			process: internal.Process{
				ProcessTypeName: sql.NullString{String: "local workflow"},
			},
			hook: Hook{
				Triggers: Trigger{
					// once config is process hook contains the name, not the code
					ProcessType: "local workflow",
				},
			},
			wantFireProcessType: true,
		},
		{
			name: "Trigger for 'reconfigure' does not fire processTypeName 'local workflow' ",
			process: internal.Process{
				ProcessTypeName: sql.NullString{String: "local workflow"},
			},
			hook: Hook{
				Triggers: Trigger{
					ProcessType: "reconfigure",
				},
			},
			wantFireProcessType: false,
		},

		// taskName alone
		{
			name: "Fires for taskName 'Test task' ",
			process: internal.Process{
				TaskName: sql.NullString{String: "Test Task"},
			},
			hook: Hook{
				Triggers: Trigger{
					TaskName: "Test Task",
				},
			},
			wantFireTaskName: true,
		},
		{
			name: "Does not fire for taskName 'Test task' trigger is for different task",
			process: internal.Process{
				TaskName: sql.NullString{String: "Test Task"},
			},
			hook: Hook{
				Triggers: Trigger{
					TaskName: "Test Task With Another name",
				},
			},
			wantFireTaskName: false,
		},

		// AccountId alone
		{
			name: "Fires for AccountId '1' ",
			process: internal.Process{
				AccountId: sql.NullInt64{Int64: 1},
			},
			hook: Hook{
				Triggers: Trigger{
					AccountId: 1,
				},
			},
			wantFireAccountId: true,
		},
		{
			name: "Does not fire for AccountId '1' trigger is looking for tenant with id '2'",
			process: internal.Process{
				AccountId: sql.NullInt64{Int64: 1},
			},
			hook: Hook{
				Triggers: Trigger{
					AccountId: 2,
				},
			},
			wantFireAccountId: false,
		},

		// CreatedBy alone
		{
			name: "Fires for created by 'Testuser' ",
			process: internal.Process{
				CreatedBy: sql.NullString{String: "Testuser"},
			},
			hook: Hook{
				Triggers: Trigger{
					CreatedBy: "Testuser",
				},
			},
			wantFireCreatedBy: true,
		},
		{
			name: "Does not fire for created by 'Testuser' trigger is looking for 'admin'",
			process: internal.Process{
				CreatedBy: sql.NullString{String: "Testuser"},
			},
			hook: Hook{
				Triggers: Trigger{
					CreatedBy: "admin",
				},
			},
			wantFireCreatedBy: false,
		},

		// Combination additive test cases from here to cover remaining function logic

		{
			name: "Should not fire, requires additional matches 1",
			process: internal.Process{
				Status: "complete",
			},
			hook: Hook{
				Triggers: Trigger{
					Status:    "complete",
					AccountId: 2,
				},
			},
			wantFireStatus: false,
		},
		{
			name: "Should not fire, requires additional matches 1",
			process: internal.Process{
				Status:          "complete",
				ProcessTypeName: sql.NullString{String: "Test Process Name"},
				AccountId:       sql.NullInt64{Int64: 2},
			},
			hook: Hook{
				Triggers: Trigger{
					Status:      "complete",
					ProcessType: "reconfigure",
					AccountId:   2,
				},
			},
			wantFireStatus: false,
		},
		{
			name: "Should not fire, requires additional matches 3",
			process: internal.Process{
				Status:          "complete",
				ProcessTypeName: sql.NullString{String: "Test Process Name"},
				TaskName:        sql.NullString{String: "Test Task Name"},
				AccountId:       sql.NullInt64{Int64: 2},
			},
			hook: Hook{
				Triggers: Trigger{
					Status:      "complete",
					ProcessType: "Test Process Name",
					TaskName:    "Different Task Name",
					AccountId:   2,
				},
			},
			wantFireStatus: false,
		},
		{
			name: "Should not fire, requires additional matches 4",
			process: internal.Process{
				Status:          "complete",
				ProcessTypeName: sql.NullString{String: "Test Process Name"},
				TaskName:        sql.NullString{String: "Test Task Name"},
				CreatedBy:       sql.NullString{String: "Test User"},
				AccountId:       sql.NullInt64{Int64: 2},
			},
			hook: Hook{
				Triggers: Trigger{
					Status:      "complete",
					ProcessType: "Test Process Name",
					TaskName:    "Test Task Name",
					CreatedBy:   "Admin",
					AccountId:   2,
				},
			},
			wantFireStatus: false,
		},
		{
			name: "Should fire, all triggers match",
			process: internal.Process{
				Status:          "complete",
				ProcessTypeName: sql.NullString{String: "Test Process Name"},
				TaskName:        sql.NullString{String: "Test Task Name"},
				CreatedBy:       sql.NullString{String: "Test User"},
				AccountId:       sql.NullInt64{Int64: 2},
			},
			hook: Hook{
				Triggers: Trigger{
					Status:      "complete",
					ProcessType: "Test Process Name",
					TaskName:    "Test Task Name",
					CreatedBy:   "Test User",
					AccountId:   2,
				},
			},
			// only want event will fire in the code
			// we use continue to move to the next hook
			wantFireStatus:      true,
			wantFireProcessType: true,
			wantFireTaskName:    true,
			wantFireCreatedBy:   true,
			wantFireAccountId:   true,
		},
	}

	// status test
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// status
			gotFireStatus := checkStatus(tc.process, tc.hook)
			if gotFireStatus != tc.wantFireStatus {
				t.Errorf("checkStatus wanted %v got %v", tc.wantFireStatus, gotFireStatus)
			}

			// processType
			gotFireProcessType := checkProcessType(tc.process, tc.hook)
			if gotFireProcessType != tc.wantFireProcessType {
				t.Errorf("checkProcessType wanted %v got %v", tc.wantFireProcessType, gotFireProcessType)
			}

			// taskName
			gotFireTaskName := checkTaskName(tc.process, tc.hook)
			if gotFireTaskName != tc.wantFireTaskName {
				t.Errorf("checkTaskName wanted %v got %v", tc.wantFireTaskName, gotFireTaskName)
			}

			// AccountId
			gotFireAccountId := checkAccountId(tc.process, tc.hook)
			if gotFireAccountId != tc.wantFireAccountId {
				t.Errorf("checkAccountId wanted %v got %v", tc.wantFireAccountId, gotFireAccountId)
			}

			// CreatedBy
			gotFireCreatedBy := checkCreatedBy(tc.process, tc.hook)
			if gotFireCreatedBy != tc.wantFireCreatedBy {
				t.Errorf("checkCreatedBy wanted %v got %v", tc.wantFireCreatedBy, gotFireCreatedBy)
			}

		})
	}

}
