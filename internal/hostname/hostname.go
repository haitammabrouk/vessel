package hostname

import (
	"fmt"
	"golang.org/x/sys/unix"
	"vessel/internal/id"
	"os"
)

func SetHostname() error {
	shortId := id.GetShortHandId(os.Args[2])
	if err := unix.Sethostname([]byte(shortId)); err != nil {
		return fmt.Errorf("set hostname: %w", err)
	}
	return nil
}