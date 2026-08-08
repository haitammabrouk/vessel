package cgroup

import (
	"context"
	"fmt"
	sdbus "github.com/coreos/go-systemd/v22/dbus"
	"github.com/godbus/dbus/v5"
	"vessel/internal/cgroup/resources"
)

// setUpUnitScopeProps set unit scope properties
func setUpUnitScopeProps(pid int, limits resources.Limits) ([]sdbus.Property, error){
	props := []sdbus.Property{
		sdbus.PropSlice("system.slice"),
		sdbus.PropPids(uint32(pid)),
		sdbus.PropDescription("vessel container scope"),
	}
	
	if limits.Memory.Max > 0 {
		props = append(props, sdbus.Property{
			Name: "MemoryMax",
			Value: dbus.MakeVariant(limits.Memory.Max),
		})
	}

	if limits.Pids.Max > 0 {
		props = append(props, sdbus.Property{
			Name: "TasksMax",
			Value: dbus.MakeVariant(limits.Pids.Max),
		})
	}
	return props, nil
}

// CreateUnitScope creates the unit scope for the container
func CreateUnitScope(containerId string, hostPid int, limits resources.Limits) error {
	ctx := context.Background()
	// connect to systemd
	conn, err := sdbus.NewSystemConnectionContext(ctx)
	if err != nil {
		return fmt.Errorf("cannot connect to systemd :%w", err)
	}
	defer conn.Close()

	unitName := fmt.Sprintf("vessel-%s.scope", containerId)
	ch := make(chan string)
	props, err := setUpUnitScopeProps(hostPid, limits)
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

func StopUnitScope(containerId string) error {
	ctx := context.Background()
	// connect to systemd
	conn, err := sdbus.NewSystemConnectionContext(ctx)
	if err != nil {
		return fmt.Errorf("cannot connect to systemd :%w", err)
	}
	defer conn.Close()
	
	unitName := fmt.Sprintf("vessel-%s.scope", containerId)
	ch := make(chan string)

	if _, err := conn.StopUnitContext(ctx, unitName, "fail", ch); err != nil {
		return fmt.Errorf("send unit stop job to systemd: %w", err)
	}

	resp := <-ch
	if resp != "done" {
		return fmt.Errorf("cannot stop unit scope for container :%w", err)
	}
	return nil
}