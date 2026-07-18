package container

import (
	"fmt"
	"os"
	"os/exec"
	"vessel/internal/cgroup"
	"golang.org/x/sys/unix"
)

func Child() error {
	if err := cgroup.SetUpCgroup(os.Getpid()); err != nil {
		return err
	}

	if err := unix.Sethostname([]byte("vessel")); err != nil {
		return fmt.Errorf("set hostname: %w", err)
	}

	cmd := exec.Command("/bin/ash")
	
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := setUpRootFs(); err != nil {
		return err
	}
	
	return cmd.Run()
}