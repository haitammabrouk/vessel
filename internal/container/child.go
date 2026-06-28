package container

import (
	"os"
	"os/exec"
)

func Child() error {
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