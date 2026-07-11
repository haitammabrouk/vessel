package cgroup

import (
	"os"
	"path/filepath"
	"fmt"
	"strconv"
)

const cgroupFs = "/sys/fs/cgroup"

func SetUpCgroup(pid int) error {
	myCgroupPath := filepath.Join(cgroupFs, "mycgroup")
	if err := os.MkdirAll(myCgroupPath, 0700); err != nil {
		return fmt.Errorf("create cgroup: %w", err)
	}

	err := addProcessToCgroup(myCgroupPath, pid)
	if err != nil {
		return err
	}

	// TODO set cgroup limits
	return nil
}

func addProcessToCgroup(myCgroupPath string, pid int) error {
	cgroupProcFilePath := filepath.Join(myCgroupPath, "cgroup.procs")
	cgroupProcFile, err := os.OpenFile(cgroupProcFilePath, os.O_WRONLY, 0)

	if err != nil {
		return fmt.Errorf("open cgroup procs file: %w", err)
	}
	defer cgroupProcFile.Close()

	_, err = cgroupProcFile.WriteString(strconv.Itoa(pid))

	if err != nil {
		return fmt.Errorf("add current process to cgroup")
	}
	return nil
}

// TODO implement
func SetUpMemoryLimits() error {
	return nil
}
// TODO implement
func SetCpuLimits() error {
	return nil
}