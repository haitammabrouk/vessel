package container

import (
	"os"
	"os/exec"
	"vessel/internal/namespace"
	"vessel/internal/cgroup"
)

func Run() error {
	cmd := exec.Command("/proc/self/exe", "child")

	cmd.SysProcAttr = namespace.SetUpNs()

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	pid := cmd.Process.Pid
	if err := cgroup.SetUpCgroup(pid); err != nil {
		return err
	}

	return cmd.Wait()
}