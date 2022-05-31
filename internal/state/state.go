package state

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

const FILE_NAME = "dozer.state"

// State holds information about the last poll against the database and
// processes which where being tracked as "executing" when the application was terminated
type State struct {
	LastPollProcessId  int       `json:"lastPollProcessId"`
	LastPollTimestamp  time.Time `json:"lastPollTimestamp"`
	ExecutingProcesses []int     `json:"executingProcesses"`
}

// HasSavedState performs a simple check to discover saved state from an application shutdown
func (*State) HasSavedState() bool {
	state := true
	if _, err := os.Stat(FILE_NAME); errors.Is(err, os.ErrNotExist) {
		state = false
	}
	return state
}

// ReadAndParse will read the state file parse the contents and make the required info available in the return
func (s *State) ReadAndParse() error {
	data, err := os.ReadFile(FILE_NAME)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, s); err != nil {
		return err
	}
	return nil
}

// CreateAndWrite will marshal state information to JSON and write it to the state file
func (s *State) CreateAndWrite() error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(FILE_NAME, data, 0644)
}

// DeleteProcessFromState removes an element from State.ExecutingProcesses state property by value
func (s *State) DeleteProcessFromState(id int) {
	var tmpState []int
	for i := range s.ExecutingProcesses {
		if s.ExecutingProcesses[i] == id {
			tmpState = append(s.ExecutingProcesses[:i], s.ExecutingProcesses[i+1:]...)
		}
	}
	s.ExecutingProcesses = tmpState
}
