package container

import (
	"os"
	"os/exec"
	"vessel/internal/hostname"
)

func Child() error {
	r := os.NewFile(3, "sync")
	buff := make([]byte, 2)
	r.Read(buff)
	r.Close()

	cmd := exec.Command("/bin/ash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := hostname.SetHostname(); err != nil {
		return err
	}

	if err := setUpRootFs(); err != nil {
		return err
	}
	
	return cmd.Run()
}