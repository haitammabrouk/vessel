package namespace

import (
	"golang.org/x/sys/unix"
)

const cloneFlags = unix.CLONE_NEWPID | 
				unix.CLONE_NEWNS

func SetUpNs() *unix.SysProcAttr {

	return &unix.SysProcAttr{
		Cloneflags: cloneFlags,
	}
}



