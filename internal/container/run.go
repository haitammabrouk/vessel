package container

import (
	"os"
	"os/exec"
	"vessel/internal/namespace"
	"vessel/internal/cgroup"
	"vessel/internal/id"
	"vessel/internal/metadata"
	"time"
)

func Run() error {
	containerId := id.GenerateRandomId()
	r, w, _ := os.Pipe()

	cmd := exec.Command("/proc/self/exe", append([]string{"child", containerId}, os.Args[2:]...)...)
	cmd.SysProcAttr = namespace.SetUpNs()
	cmd.ExtraFiles = []*os.File{r}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Start(); err != nil {
		return err
	}

	r.Close()

	hostPid := cmd.Process.Pid
	if err := cgroup.CreateUnitScope(containerId, hostPid); err != nil {
		return err
	}
	
	if err := metadata.InitContainerConfig(containerId, time.Now().Format(time.RFC3339Nano)) ; err != nil {
			return err
	}

	w.Write([]byte("go"))
	w.Close()

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}