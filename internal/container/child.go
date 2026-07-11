package container

import (
	"os"
	"os/exec"
	"vessel/internal/cgroup"
)

func Child() error {
	if err := cgroup.SetUpCgroup(os.Getpid()); err != nil {
		return err
	}

	if err := setUpRootFs(); err != nil {
		return err
	}
	
	cmd := exec.Command("/bin/ash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}