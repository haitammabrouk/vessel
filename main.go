package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"golang.org/x/sys/unix"
	"path/filepath"
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

	// swap root fs
	pivotRoot()

	// mount the proc and associate it with the pid table of the new PID namespace
	if err := unix.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		fmt.Fprintf(os.Stderr, "Error mounting proc: %v\n", err)
		os.Exit(1)
	}
	
	cmd := exec.Command("/bin/ash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting the bash process: %v\n", err)
		os.Exit(1)
	}
	
	if err := cmd.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "Error waiting for the bash process: %v\n", err)
		os.Exit(1)
	}
}

func run() {
	// rerun the main program with a child argument
	cmd := exec.Command("/proc/self/exe", "child")
	// set up namespaces (PID, mount, and user)
	cmd.SysProcAttr = &unix.SysProcAttr{
		Cloneflags: unix.CLONE_NEWUSER | unix.CLONE_NEWPID | unix.CLONE_NEWNS,

		UidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getuid(), Size: 1},
		},
		GidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getgid(), Size: 1},
		},

		GidMappingsEnableSetgroups: false,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running the child process: %v\n", err)
		os.Exit(1)
	}
}

func pivotRoot() error {
	newRootFsPath, err := filepath.Abs("./rootfs")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting abs path of new rootfs: %v\n", err)
		os.Exit(1)
	}

	// make mount points private so that they don't propagate to the host
	if err := unix.Mount("", "/", "", unix.MS_PRIVATE|unix.MS_REC, ""); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting mount propagation to private: %v\n", err)
		os.Exit(1)
	}

	// making the new root fs mountable
	if err := unix.Mount(newRootFsPath, newRootFsPath, "", unix.MS_BIND|unix.MS_REC, ""); err != nil {
		fmt.Fprintf(os.Stderr, "Error making new rootfs mountable: %v\n", err)
		os.Exit(1)
	}

	oldRootFsPath := filepath.Join(newRootFsPath, "old_root")
	// create the folder inside the new rootfs that will hold old rootfs
	if err := os.MkdirAll(oldRootFsPath, 0700); err != nil {
		fmt.Fprintf(os.Stderr, "Errot creating old root folder inside new rootfs: %v\n", err)
		os.Exit(1)
	}

	// swap the old rootfs with the new rootfs
	if err := unix.PivotRoot(newRootFsPath, oldRootFsPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error pivoting rootfs: %v\n", err)
		os.Exit(1)
	}

	// ensure cwd set to new root
	if err := unix.Chdir("/"); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting the cwd to root: %v\n", err)
		os.Exit(1)
	}

	// umount old rootfs
	if err := unix.Unmount("./old_root", unix.MNT_DETACH); err != nil {
		fmt.Fprintf(os.Stderr, "Error umounting old rootfs: %v\n", err)
		os.Exit(1)
	}

	// remove old rootfs
	if err := os.RemoveAll("./old_root"); err != nil {
		fmt.Fprintf(os.Stderr, "Error removing old rootfs: %v\n", err)
		os.Exit(1)
	}

	return nil
}