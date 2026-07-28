package hostname

import (
	"fmt"
	"golang.org/x/sys/unix"
	"vessel/internal/id"
)

func SetHostname() error {
	id := id.GetShortHandId(id.GenerateRandomId())

	if err := unix.Sethostname([]byte(id)); err != nil {
		return fmt.Errorf("set hostname: %w", err)
	}
	return nil
}