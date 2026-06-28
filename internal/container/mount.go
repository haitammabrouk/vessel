package container

import (
	"fmt"
	"golang.org/x/sys/unix"
)

func makeMountsPrivate() error {
	if err := unix.Mount("", "/", "", unix.MS_PRIVATE | unix.MS_REC, ""); err != nil {
		return fmt.Errorf("make mounts private: %w", err)
	}
	return nil
}

func mountProcFs() error {
	if err := unix.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		return fmt.Errorf("mount procfs: %w", err)
	}
	return nil
}