package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"golang.org/x/sys/unix"
	"path/filepath"
)

func main() {
    switch os.Args[1] {
    case "run":
        if err := run(); err != nil {
			fmt.Fprintf(os.Stderr, "setup namespaces failed: %v\n", err)
			os.Exit(1)
		}
    case "child":
        if err := child(); err != nil {
			fmt.Fprintf(os.Stderr, "init container failed: %v\n", err)
			os.Exit(1)
		}
    }
}

func child() error {
	// swap root fs
	if err := setUpRootFs(); err != nil {
		return err
	}
	
	cmd := exec.Command("/bin/ash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

func run() error {
	cmd := exec.Command("/proc/self/exe", "child")

	cmd.SysProcAttr = &unix.SysProcAttr{
		Cloneflags: unix.CLONE_NEWUSER | unix.CLONE_NEWPID | unix.CLONE_NEWNS,

		UidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getuid(), Size: 1},
		},
		GidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getgid(), Size: 1},
		},
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func makeMountsPrivate() error {
	if err := unix.Mount("", "/", "", unix.MS_PRIVATE | unix.MS_REC, ""); err != nil {
		return fmt.Errorf("make mounts private: %w", err)
	}
	return nil
}

func detachOldRootFs() error {
	if err := unix.Unmount("./old_root", unix.MNT_DETACH); err != nil {
		return fmt.Errorf("umount old rootfs: %w", err)
	}

	if err := os.RemoveAll("./old_root"); err != nil {
		return fmt.Errorf("remove old rootfs: %w", err)
	}
	return nil
}

func mountProcFs() error {
	if err := unix.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		return fmt.Errorf("mount procfs: %w", err)
	}
	return nil
}

func bindNewMountRootFs(newRootFsPath string) error {
	if err := unix.Mount(newRootFsPath, newRootFsPath, "", unix.MS_BIND | unix.MS_REC, ""); err != nil {
		return fmt.Errorf("bind new rootfs: %w", err)
	}
	return nil
}

func createOldRootFsPlaceholder(oldRootFsPath string) error {
	if err := os.MkdirAll(oldRootFsPath, 0700); err != nil {
		return fmt.Errorf("create old rootfs placeholder: %w", err)
	}
	return nil
}

func setWorkingDirToRoot() error {
	if err := unix.Chdir("/"); err != nil {
		return fmt.Errorf("change working directory to root: %w", err)
	}
	return nil
}

func setUpRootFs() error {
	newRootFsPath, err := filepath.Abs("./rootfs")

	if err != nil {
		return err
	}

	oldRootFsPath := filepath.Join(newRootFsPath, "old_root")

	if err := makeMountsPrivate(); err != nil {
		return err
	}

	if err := bindNewMountRootFs(newRootFsPath); err != nil {
		return err
	}

	if err := createOldRootFsPlaceholder(oldRootFsPath); err != nil {
		return err
	}

	if err := unix.PivotRoot(newRootFsPath, oldRootFsPath); err != nil {
		return err
	}

	if err := setWorkingDirToRoot(); err != nil {
		return err
	}

	if err := mountProcFs(); err != nil {
		return err
	}

	if err := detachOldRootFs(); err != nil {
		return err
	}

	return nil
}