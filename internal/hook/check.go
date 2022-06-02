package hook

import (
	"context"
	"fmt"

	"github.com/spoonboy-io/dozer/internal"
	"github.com/spoonboy-io/koan"
)

// TODO we need he logger

// CheckProcess will check a process against the configuration to determine if
// it is an event that should trigger a call webhook
func CheckProcess(process *internal.Process, logger *koan.Logger, ctx context.Context) {
	// go through all the hook config
	for i := range config {
		// config uses code we need to search on name, we swap the code for the name in the config
		processName, err := getProcessTypeName(config[i].Triggers.ProcessType)
		if err == nil {
			config[i].Triggers.ProcessType = processName
		}

		if fire := checkStatus(process, &config[i].Hook); fire {
			if err := fireWebhook(process, &config[i].Hook, logger, ctx); err != nil {
				warnMsg := fmt.Sprintf("Failed to fire webhook on status (hook: '%s', url: '%s', process id: '%d')",
					config[i].Hook.Description, config[i].Hook.URL, process.Id)
				logger.Warn(warnMsg)
			}
			continue
		}

		if fire := checkProcessType(process, &config[i].Hook); fire {
			if err := fireWebhook(process, &config[i].Hook, logger, ctx); err != nil {
				warnMsg := fmt.Sprintf("Failed to fire webhook on procesType (hook: '%s', url: '%s', process id: '%d')",
					config[i].Hook.Description, config[i].Hook.URL, process.Id)
				logger.Warn(warnMsg)
			}
			continue
		}

		if fire := checkTaskName(process, &config[i].Hook); fire {
			if err := fireWebhook(process, &config[i].Hook, logger, ctx); err != nil {
				warnMsg := fmt.Sprintf("Failed to fire webhook on taskName (hook: '%s', url: '%s', process id: '%d')",
					config[i].Hook.Description, config[i].Hook.URL, process.Id)
				logger.Warn(warnMsg)
			}
			continue
		}

		if fire := checkAccountId(process, &config[i].Hook); fire {
			if err := fireWebhook(process, &config[i].Hook, logger, ctx); err != nil {
				warnMsg := fmt.Sprintf("Failed to fire webhook on accountId (hook: '%s', url: '%s', process id: '%d')",
					config[i].Hook.Description, config[i].Hook.URL, process.Id)
				logger.Warn(warnMsg)
			}
			continue
		}

		if fire := checkCreatedBy(process, &config[i].Hook); fire {
			if err := fireWebhook(process, &config[i].Hook, logger, ctx); err != nil {
				warnMsg := fmt.Sprintf("Failed to fire webhook on createdBy (hook: '%s', url: '%s', process id: '%d')",
					config[i].Hook.Description, config[i].Hook.URL, process.Id)
				logger.Warn(warnMsg)
			}
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
