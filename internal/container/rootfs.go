package container

import (
	"fmt"
	"os"
	"golang.org/x/sys/unix"
	"path/filepath"
)

func bindMountRootFs(newRootFsPath string) error {
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

func detachOldRootFs() error {
	if err := unix.Unmount("./old_root", unix.MNT_DETACH); err != nil {
		return fmt.Errorf("umount old rootfs: %w", err)
	}

	if err := os.RemoveAll("./old_root"); err != nil {
		return fmt.Errorf("remove old rootfs: %w", err)
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

	if err := bindMountRootFs(newRootFsPath); err != nil {
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