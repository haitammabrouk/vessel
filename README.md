# Vessel

Vessel is a lightweight container written in Go for learning how Linux containers work under the hood.

The project aims to build a minimal container from scratch without relying on existing container engines (no runc, no libcontainer). It focuses on understanding Linux primitives such as:

- Linux namespaces (PID, mount, UTS, IPC)
- Root filesystems (`pivot_root`-based rootfs switch)
- Cgroups v2, managed as real **systemd scopes** (`vessel-<id>.scope` under
  `system.slice`, created via systemd's D-Bus API) rather than raw `mkdir`
  under `/sys/fs/cgroup` — see [TODO.md](./TODO.md) §3 for why
- Container identity and metadata (random 256-bit container ID, JSON config
  under `/var/lib/vessel/containers/<id>/`)

> **Note:** This is an educational project and is **not** intended for production use.

## Prerequisites

- Linux, booted with **systemd** as PID 1, with a reachable system D-Bus —
  container cgroups are created as transient systemd scopes, so this is a
  hard requirement, not just a recommendation
- Go 1.25+ (or your installed Go version)
- Run as root (or via `sudo`) — namespace/cgroup setup needs the relevant privileges

## Getting Started

### 1. Download the Alpine miniroot filesystem

```bash
./scripts/pull-rootfs.sh
```

This downloads the Alpine Mini RootFS and extracts it into the `rootfs/` directory.

### 2. Build the project

```bash
go build -o bin/vessel ./cmd/vessel
```

### 3. Run the container

```bash
sudo ./bin/vessel run
```

If everything is configured correctly, you'll be dropped into an Alpine shell running inside its own namespaces.

Optionally cap memory:

```bash
sudo ./bin/vessel run --memory 256m
```

### 4. Inspect the container from the host

Each run gets a random ID, visible as the container's hostname from inside
the shell (`hostname`). From the host, while the container is running:

```bash
# real systemd scope for this container, with live resource accounting
systemctl status vessel-<id>.scope

# on-disk metadata for this container
cat /var/lib/vessel/containers/<id>/config.json
```

## License

This project is provided for educational purposes.