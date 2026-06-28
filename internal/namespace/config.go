package namespace

import (
	"os"
	"syscall"
	"golang.org/x/sys/unix"
)

const cloneFlags = unix.CLONE_NEWUSER | 
				unix.CLONE_NEWPID | 
				unix.CLONE_NEWNS

func SetUpNs() *unix.SysProcAttr {
	hostUID := os.Getuid()
	hostGID := os.Getgid()

	return &unix.SysProcAttr{
		Cloneflags: cloneFlags,

		UidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: hostUID, Size: 1},
		},
		GidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: hostGID, Size: 1},
		},
	}
}



