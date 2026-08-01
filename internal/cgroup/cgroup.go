package cgroup

import (
	"context"
	"fmt"
	"vessel/internal/cli"
	sdbus "github.com/coreos/go-systemd/v22/dbus"
	"github.com/godbus/dbus/v5"
)

// setUpUnitScopeProps set unit scope properties
func setUpUnitScopeProps(pid int) ([]sdbus.Property, error){
	limits, err := cli.ParseOptions()
	if err != nil {
		return nil, err
	}
	return []sdbus.Property{
		sdbus.PropSlice("system.slice"),
		sdbus.PropPids(uint32(pid)),
		sdbus.PropDescription("vessel container scope"),
		{
			Name: "MemoryMax",
			Value: dbus.MakeVariant(uint64(limits.Memory.Max)),
		},
		{
			Name: "MemorySwapMax",
			Value: dbus.MakeVariant(uint64(limits.Memory.SwapMax)),
		},
	}, nil
}

// CreateUnitScope creates the unit scope for the container
func CreateUnitScope(containerId string, hostPid int) error {
	ctx := context.Background()
	// connect to systemd
	conn, err := sdbus.NewSystemConnectionContext(ctx)
	if err != nil {
		return fmt.Errorf("cannot connect to systemd :%w", err)
	}
	defer conn.Close()

	unitName := fmt.Sprintf("vessel-%s.scope", containerId)
	ch := make(chan string)
	props, err := setUpUnitScopeProps(hostPid)
	if err != nil {
		return err
	}

	if _, err := conn.StartTransientUnitContext(ctx, unitName, "fail", props, ch); err != nil {
		return fmt.Errorf("send unit creation job to systemd :%w", err)
	}

	resp := <-ch
	if resp != "done" {
		return fmt.Errorf("cannot create unit scope for container :%w", err)
	}
	return nil
}