package cgroup

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"vessel/internal/cgroup/resources"
	"vessel/internal/cli"
	"golang.org/x/sys/unix"
)

const cgroupFs = "/sys/fs/cgroup"
var cgroupPath string

func SetUpCgroup(pid int) error {
	cgroupPath = filepath.Join(cgroupFs, "mycgroup")
	if err := os.MkdirAll(cgroupPath, 0700); err != nil {
		return fmt.Errorf("create cgroup: %w", err)
	}

	err := addProcessToCgroup(pid)
	if err != nil {
		return err
	}

	resourceLimits, err := cli.ParseOptions()
	if err != nil {
		return err
	}

	if err = setMemoryLimits(resourceLimits.Memory); err != nil {
		return err
	}

	if err = unix.Unshare(unix.CLONE_NEWCGROUP); err != nil {
		return fmt.Errorf("unshare cgroup: %w", err)
	}

	return nil
}

func addProcessToCgroup(pid int) error {
	cgroupProcFilePath := filepath.Join(cgroupPath, "cgroup.procs")
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

func setMemoryLimits(memoryLimits resources.Memory) error {
	memoryMaxFilePath := filepath.Join(cgroupPath, "memory.max")
	memoryMaxFile, err := os.OpenFile(memoryMaxFilePath, os.O_WRONLY, 0)

	if err != nil {
		return fmt.Errorf("open memory.max file: %w", err)
	}
	defer memoryMaxFile.Close()

	if memoryLimits.Max == 0 {
		_, err = memoryMaxFile.WriteString("max")
	} else {
		_, err = memoryMaxFile.WriteString(strconv.FormatInt(memoryLimits.Max, 10))
	}

	if err != nil {
		return fmt.Errorf("set memory max limit")
	}

	memorySwapMaxFilePath := filepath.Join(cgroupPath, "memory.swap.max")
	memorySwapMaxFile, err := os.OpenFile(memorySwapMaxFilePath, os.O_WRONLY, 0)

	if err != nil {
		return fmt.Errorf("open memory.swap.max file: %w", err)
	}
	defer memorySwapMaxFile.Close()

	if memoryLimits.SwapMax == 0 {
		_, err = memorySwapMaxFile.WriteString("max")
	} else {
		_, err = memorySwapMaxFile.WriteString(strconv.FormatInt(memoryLimits.SwapMax, 10))
	}

	if err != nil {
		return fmt.Errorf("set memory swap max limit")
	}

	return nil
}