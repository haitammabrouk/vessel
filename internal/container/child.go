package container

import (
	"os"
	"os/exec"
	"vessel/internal/cgroup"
	"vessel/internal/hostname"
	"vessel/internal/id"
	"vessel/internal/metadata"
	"time"
)

func Child() error {
	containerId := id.GenerateRandomId()

	if err := cgroup.SetUpCgroup(os.Getpid()); err != nil {
		return err
	}

	if err := hostname.SetHostname(containerId); err != nil {
		return err
	}
	
	if err := metadata.InitContainerConfig(containerId, time.Now().Format(time.RFC3339Nano)) ; err != nil {
			return err
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