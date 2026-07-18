# TODO

Working today: PID + mount namespaces + uts namespaces, cgroup v2 memory limits, pivot_root-based
rootfs switch, hardcoded `/bin/ash` launch.

## 1. Namespace isolation (currently only PID + mount)
- [ ] Add `CLONE_NEWIPC` — isolate SysV IPC / POSIX message queues
- [ ] Add `CLONE_NEWNET` — container currently shares the host network stack entirely
- [ ] Add `CLONE_NEWUSER` + UID/GID mapping (`/proc/<pid>/uid_map`, `gid_map`) — needed
      for rootless operation

## 2. Networking (currently none)
- [ ] Create a veth pair, move one end into the container's net namespace
- [ ] Bring up `lo` inside the container (namespaces start with loopback down)
- [ ] Bridge/NAT on the host side + iptables MASQUERADE rule for outbound traffic
- [ ] Basic DNS (`/etc/resolv.conf` in the mounted rootfs)

## 3. Cgroups — fix existing gaps before adding more controllers
- [ ] `cgroupPath` is hardcoded to `mycgroup` (`internal/cgroup/cgroup.go:17`) —
      concurrent containers collide on the same cgroup; generate a unique ID per run
- [ ] Nothing removes the cgroup directory on exit — dirs leak on every run
- [ ] `cgroup.SetUpCgroup` calls `cli.ParseOptions()` itself instead of receiving
      limits as a parameter (`internal/cgroup/cgroup.go:27`) — reversed dependency,
      re-parses `os.Args` a second time in the child process; pass
      `resources.ResouceLimits` in from `Child()` instead
- [ ] Add CPU limits (`cpu.max`) and `pids.max` (fork-bomb protection)
- [ ] Add cgroup v1 fallback/detection, or fail with a clear error on v1-only hosts

## 4. Process lifecycle
- [ ] Reap zombies — the shell becomes PID 1 in the new PID namespace with no
      reaping logic; orphaned children will accumulate as zombies
- [ ] Forward signals from the `run` process to the `child`/container process
- [ ] Rollback/cleanup on partial failure (e.g. cgroup created but rootfs setup fails)

## 5. Rootfs / mounts
- [ ] `./rootfs` path is hardcoded (`internal/container/rootfs.go:43`) — no way to
      point at a different image/rootfs per run
- [ ] Command to run is hardcoded to `/bin/ash` (`internal/container/child.go:14`) —
      should come from CLI args
- [ ] Unmount/cleanup the bind-mounted rootfs after the container exits (only
      `old_root` is detached today, the outer bind mount is never undone)
- [ ] Read-only rootfs support + writable overlay (overlayfs) so the base image
      isn't mutated between runs

## 6. Security hardening (none of this exists yet)
- [ ] Drop Linux capabilities (keep only what's needed, drop `CAP_SYS_ADMIN` etc.
      after setup is done)
- [ ] Set `no_new_privs`
- [ ] Seccomp filtering (even a basic default profile)
- [ ] Mount `/proc`, `/sys` with restrictive flags (`nosuid`, `noexec`, `nodev`
      where sensible)

## 7. CLI/UX
- [ ] `os.Args[1]` indexing in `main.go` panics if no args are given — add usage
      output and bounds-checked parsing
- [ ] Container identity: name/ID generation so runs can be tracked, listed, stopped
- [ ] `vessel ps` / `vessel stop <id>` — minimal state file (PID + cgroup path) for
      lifecycle awareness beyond a single foreground run
- [ ] `vessel exec <id> <cmd>` — join an existing container's namespaces via `setns`
