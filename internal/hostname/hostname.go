package hostname

import (
	"fmt"
	"golang.org/x/sys/unix"
	"vessel/internal/id"
)

func SetHostname(containerId string) error {
	shortId := id.GetShortHandId(containerId)
	if err := unix.Sethostname([]byte(shortId)); err != nil {
		return fmt.Errorf("set hostname: %w", err)
	}
	return nil
}