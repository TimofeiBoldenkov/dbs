package raminfoprovider

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type RAMInfo struct {
	TotalRAMInKB	uint64	`json:"total_ram_kb"`
	UsedRAMInKB		uint64 	`json:"used_ram_kb"` // UsedRAMInKB = MemTotal - MemAvailable
}

type RAMInfoProvider struct {}

// returns RamInfo
func (RAMInfoProvider) GetInfo() (any, error) {
	var ramInfo RAMInfo

	filepath := "/proc/meminfo"
	ramInfoFile, err := os.Open(filepath)
	if err != nil {
		return RAMInfo{}, fmt.Errorf("can't open RAM info file: %v", err.Error())
	}
	defer ramInfoFile.Close()

	scanner := bufio.NewScanner(ramInfoFile)

	totalRAM, err := scanUintFromFirstLine(scanner)
	if err != nil {
		return RAMInfo{}, fmt.Errorf("can't parse %v: %v", filepath, err.Error())
	}
	ramInfo.TotalRAMInKB = totalRAM

	scanner.Scan() // skip second line

	availableRAM, err := scanUintFromFirstLine(scanner)
	if err != nil {
		return RAMInfo{}, fmt.Errorf("can't parse %v: %v", filepath, err.Error())
	}
	ramInfo.UsedRAMInKB = ramInfo.TotalRAMInKB - availableRAM

	return ramInfo, nil
}

func scanUintFromFirstLine(scanner *bufio.Scanner) (uint64, error) {
	if scanner.Scan() {
		number, ok := getFirstNumber(scanner.Text())
		if !ok {
			return 0, fmt.Errorf("can't parse uint64 from line: %v", scanner.Text())
		}
		return number, nil
	} else {
		if err := scanner.Err(); err != nil {
			return 0, fmt.Errorf("can't scan: %v", err.Error())
		} else {
			return 0, errors.New("cant scan: EOF")
		}
	}
}

func getFirstNumber(str string) (uint64, bool) {
	re := regexp.MustCompile(`\d+`)

	numberStr := re.FindString(str)
	if numberStr != "" {
		number, err := strconv.ParseUint(numberStr, 10, 64)
		if err != nil {
			return 0, false
		}

		return number, true
	} else {
		return 0, false
	}
}
