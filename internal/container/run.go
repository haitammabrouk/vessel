package container

import (
	"os"
	"os/exec"
	"vessel/internal/namespace"
)

func Run() error {
	cmd := exec.Command("/proc/self/exe", "child")

	cmd.SysProcAttr = namespace.SetUpNs()

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}