package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"golang.org/x/sys/unix"
)

func main() {
    switch os.Args[1] {
    case "run":
        run()
    case "child":
        child()
    }
}

func child() {
	// start the child process inside the newly created namespaces
	if os.Args[1] == "child" {
		cmd := exec.Command("/bin/bash", "--noprofile", "--norc")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// mount the proc and associate it with the pid table of the new PID namespace
		if err := unix.Mount("proc", "/proc", "proc", 0, ""); err != nil {
			fmt.Fprintf(os.Stderr, "Error mounting proc: %v\n", err)
			os.Exit(1)
		}

		if err := cmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting the bash process: %v\n", err)
			os.Exit(1)
		}

		if err := cmd.Wait(); err != nil {
			fmt.Fprintf(os.Stderr, "Error waiting for the bash process: %v\n", err)
			os.Exit(1)
		}
	}
}

func run() {
	// rerun the main program with a child argument
	cmd := exec.Command("/proc/self/exe", "child")
	// set up namespaces (PID, mount, and user)
	cmd.SysProcAttr = &unix.SysProcAttr{
		Cloneflags: unix.CLONE_NEWPID | unix.CLONE_NEWNS | unix.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{ContainerID : 0, HostID : os.Getuid(), Size : 1},
		},
		GidMappings: []syscall.SysProcIDMap{
			{ContainerID : 0, HostID : os.Getgid(), Size : 1},
		},
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running the child process: %v\n", err)
		os.Exit(1)
	}
}