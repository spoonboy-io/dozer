package hook

import (
	"fmt"

	"github.com/spoonboy-io/dozer/internal"
)

// CheckProcess will check a process against the configuration to determine if
// it is an event that should trigger a call webhook
func CheckProcess(process *internal.Process) {
	// go through all the hook config
	for i := range config {
		// config uses code we need to search on name, we swap the code for the name in the config
		processName, err := getProcessTypeName(config[i].Triggers.ProcessType)
		if err == nil {
			config[i].Triggers.ProcessType = processName
		}

		if fire := checkStatus(process, &config[i].Hook); fire {
			fmt.Println("firing on status", process.Id, process.Status)
			continue
		}

		if fire := checkProcessType(process, &config[i].Hook); fire {
			fmt.Println("firing on process type", process.Id, process.ProcessTypeName.String)
			continue
		}

		if fire := checkTaskName(process, &config[i].Hook); fire {
			fmt.Println("firing on task name", process.Id, process.TaskName.String)
			continue
		}

		if fire := checkAccountId(process, &config[i].Hook); fire {
			fmt.Println("firing on account id", process.Id, process.AccountId.Int64)
			continue
		}

		if fire := checkCreatedBy(process, &config[i].Hook); fire {
			fmt.Println("firing on created by", process.Id, process.CreatedBy.String)
			continue
		}
	}
}

func checkStatus(process *internal.Process, hook *Hook) bool {
	var fire bool
	if process.Status == "" {
		return fire
	}
	if process.Status == "executing" {
		// we should never be in here as processes with executing status
		// should not be passed for inspection, we monitor until done/failed
		return fire
	}
	if hook.Triggers.Status == process.Status {
		fire = true
		if hook.Triggers.ProcessType != "" {
			fire = false
			if hook.Triggers.ProcessType == process.ProcessTypeName.String {
				fire = true
			} else {
				return fire
			}
		}
		if hook.Triggers.TaskName != "" {
			fire = false
			if hook.Triggers.TaskName == process.TaskName.String {
				fire = true
			} else {
				return fire
			}
		}

		if hook.Triggers.AccountId != 0 {
			fire = false
			if hook.Triggers.AccountId == int(process.AccountId.Int64) {
				fire = true
			} else {
				return fire
			}
		}

		if hook.Triggers.CreatedBy != "" {
			fire = false
			if hook.Triggers.CreatedBy == process.CreatedBy.String {
				fire = true
			} else {
				return fire
			}
		}
	}

	return fire
}

func checkProcessType(process *internal.Process, hook *Hook) bool {
	var fire bool
	if process.ProcessTypeName.String == "" {
		return fire
	}
	if hook.Triggers.ProcessType == process.ProcessTypeName.String {
		fire = true
		if hook.Triggers.Status != "" {
			fire = false
			if hook.Triggers.Status == process.Status {
				fire = true
			} else {
				return fire
			}
		}
		if hook.Triggers.TaskName != "" {
			fire = false
			if hook.Triggers.TaskName == process.TaskName.String {
				fire = true
			} else {
				return fire
			}
		}

		if hook.Triggers.AccountId != 0 {
			fire = false
			if hook.Triggers.AccountId == int(process.AccountId.Int64) {
				fire = true
			} else {
				return fire
			}
		}

		if hook.Triggers.CreatedBy != "" {
			fire = false
			if hook.Triggers.CreatedBy == process.CreatedBy.String {
				fire = true
			} else {
				return fire
			}
		}
	}

	return fire
}

func checkTaskName(process *internal.Process, hook *Hook) bool {
	var fire bool
	if process.TaskName.String == "" {
		return fire
	}
	if hook.Triggers.TaskName == process.TaskName.String {
		fire = true
		if hook.Triggers.Status != "" {
			fire = false
			if hook.Triggers.Status == process.Status {
				fire = true
			} else {
				return fire
			}
		}
		if hook.Triggers.ProcessType != "" {
			fire = false
			if hook.Triggers.ProcessType == process.ProcessTypeName.String {
				fire = true
			} else {
				return fire
			}
		}

		if hook.Triggers.AccountId != 0 {
			fire = false
			if hook.Triggers.AccountId == int(process.AccountId.Int64) {
				fire = true
			} else {
				return fire
			}
		}

		if hook.Triggers.CreatedBy != "" {
			fire = false
			if hook.Triggers.CreatedBy == process.CreatedBy.String {
				fire = true
			} else {
				return fire
			}
		}
	}

	return fire
}

func checkAccountId(process *internal.Process, hook *Hook) bool {
	var fire bool
	if int(process.AccountId.Int64) == 0 {
		return fire
	}
	if hook.Triggers.AccountId == int(process.AccountId.Int64) {
		fire = true
		if hook.Triggers.ProcessType != "" {
			fire = false
			if hook.Triggers.ProcessType == process.ProcessTypeName.String {
				fire = true
			} else {
				return fire
			}
		}
		if hook.Triggers.TaskName != "" {
			fire = false
			if hook.Triggers.TaskName == process.TaskName.String {
				fire = true
			} else {
				return fire
			}
		}

		if hook.Triggers.Status != "" {
			fire = false
			if hook.Triggers.Status == process.Status {
				fire = true
			} else {
				return fire
			}
		}

		if hook.Triggers.CreatedBy != "" {
			fire = false
			if hook.Triggers.CreatedBy == process.CreatedBy.String {
				fire = true
			} else {
				return fire
			}
		}
	}

	return fire
}

func checkCreatedBy(process *internal.Process, hook *Hook) bool {
	var fire bool
	if process.CreatedBy.String == "" {
		return fire
	}
	if hook.Triggers.CreatedBy == process.CreatedBy.String {
		fire = true
		if hook.Triggers.ProcessType != "" {
			fire = false
			if hook.Triggers.ProcessType == process.ProcessTypeName.String {
				fire = true
			} else {
				return fire
			}
		}
		if hook.Triggers.TaskName != "" {
			fire = false
			if hook.Triggers.TaskName == process.TaskName.String {
				fire = true
			} else {
				return fire
			}
		}

		if hook.Triggers.AccountId != 0 {
			fire = false
			if hook.Triggers.AccountId == int(process.AccountId.Int64) {
				fire = true
			} else {
				return fire
			}
		}

		if hook.Triggers.Status != "" {
			fire = false
			if hook.Triggers.Status == process.Status {
				fire = true
			} else {
				return fire
			}
		}
	}

	return fire
}
