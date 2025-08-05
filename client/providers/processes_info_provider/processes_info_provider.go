package processesinfoprovider

import (
	"fmt"
	"os"
	"strconv"
)

type ProcessesInfoProvider struct {}

type ProcessInfo struct {
	PID	uint16	`json:"pid"`
	ExePath	string	`json:"exe"`
}

type ProcessesInfo struct {
	Processes []ProcessInfo
}

func (ProcessesInfoProvider) GetInfo() (any, error) {
	const PROCESSES_DIR = "/proc"

	var processesInfo ProcessesInfo

	processes, err := os.ReadDir(PROCESSES_DIR)
	if err != nil {
		return nil, fmt.Errorf("can't read %v", PROCESSES_DIR)
	}

	for _, dirEntry := range processes {
		if !dirEntry.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(dirEntry.Name())
		if err != nil || !(pid >= 1 && pid <= 32767) {
			continue
		}

		exePath, err := os.Readlink(PROCESSES_DIR + "/" + dirEntry.Name() + "/exe")
		if err != nil {
			continue
		}

		processesInfo.Processes = append(processesInfo.Processes, ProcessInfo{uint16(pid), exePath})
	}

	return processesInfo, nil
}
