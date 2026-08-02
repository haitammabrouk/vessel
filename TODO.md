# TODO

Working today: PID + mount + UTS + IPC namespaces, pivot_root-based rootfs
switch, hardcoded `/bin/ash` launch, random 256-bit container ID (used as
hostname), on-disk container metadata (`/var/lib/vessel/containers/<id>/config.json`),
cgroup v2 memory limits — no longer a raw `mkdir` under `/sys/fs/cgroup`,
container cgroups are now real systemd scopes (`vessel-<id>.scope` under
`system.slice`, created via D-Bus `StartTransientUnit`), and parent/child
startup is synchronized over a pipe so the child doesn't fork the shell
before it's actually inside the scope's cgroup.

**Requires systemd as PID 1 with a reachable system D-Bus** — this is a new
hard dependency introduced by the systemd cgroup migration; vessel is no
longer usable on non-systemd init systems (see §3).

## 1. Namespace isolation (currently PID + mount + UTS + IPC)
- [ ] Add `CLONE_NEWNET` — container currently shares the host network stack entirely
- [ ] Add `CLONE_NEWUSER` + UID/GID mapping (`/proc/<pid>/uid_map`, `gid_map`) — needed
      for rootless operation

## 2. Networking (currently none)
- [ ] Create a veth pair, move one end into the container's net namespace
- [ ] Bring up `lo` inside the container (namespaces start with loopback down)
- [ ] Bridge/NAT on the host side + iptables MASQUERADE rule for outbound traffic
- [ ] Basic DNS (`/etc/resolv.conf` in the mounted rootfs)

## 3. Cgroups (systemd-managed scopes)
- [x] ~~`cgroupPath` hardcoded to `mycgroup`, concurrent containers collide~~ —
      replaced entirely by systemd transient scopes named `vessel-<id>.scope`,
      one per container, no shared path
- [x] ~~Nothing removes the cgroup directory on exit~~ — systemd owns the
      scope now, so it GCs the cgroup automatically when the last process
      exits (`cgroup.events` `populated` flag → 0). **Worth a one-time manual
      check** (`systemctl status vessel-<id>.scope` after a run exits, confirm
      `could not be found`) rather than fully trusting it unverified.
- [x] ~~`setUpUnitScopeProps` calls `cli.ParseOptions()` itself instead of
      receiving `resources.Limits` as a parameter~~ — `Run()` now parses
      flags once and passes `resources.Limits` down to `CreateUnitScope`;
      `internal/cgroup` no longer imports `cli` at all
- [ ] Add CPU limits (`CPUQuota=` property) and `TasksMax=` (fork-bomb
      protection) — same D-Bus property pattern as the existing `MemoryMax`
- [ ] Add cgroup v1 fallback/detection, or fail with a clear error on v1-only
      hosts — lower priority now that systemd brokers the cgroup itself, but
      `mountCgroup2` (`internal/container/mount.go:22`) still hardcodes a
      cgroup2 mount inside the container regardless of host cgroup version

## 4. Process lifecycle & metadata
- [x] ~~Store metadata about the running container (StartedAt, Id etc)~~ —
      `internal/metadata/metadata.go` writes `config.json` (ID, Created,
      ConfigPath) to `/var/lib/vessel/containers/<id>/` on every `run`
- [ ] Nothing removes the metadata directory on exit either — same leak shape
      the old cgroup dirs had; `vessel stop`/exit should clean up
      `/var/lib/vessel/containers/<id>/` (or explicitly keep it for `vessel ps`
      history — decide which, then implement)
- [ ] Write hostname inside a hostname file in the metadata dir of the running container
- [x] ~~`Child()`/`hostname.SetHostname()` recover the container ID by
      independently re-reading `os.Args[2]`~~ — `main.go` now passes
      `containerId` into `container.Child(containerId)` explicitly, which
      passes it straight through to `hostname.SetHostname(containerId)`
- [ ] Sync pipe read/write in `Run()`/`Child()` ignores errors on
      `r.Read(buff)` / `w.Write(...)` (`internal/container/child.go:12`,
      `internal/container/run.go:39`) — fine in practice since sub-`PIPE_BUF`
      writes are atomic on Linux, but worth handling explicitly for a
      partial-read/failed-write edge case
- [ ] Reap zombies — the shell becomes PID 1 in the new PID namespace with no
      reaping logic; orphaned children will accumulate as zombies
- [ ] Forward signals from the `run` process to the `child`/container process
- [ ] Rollback/cleanup on partial failure (e.g. scope created but rootfs setup fails)

## 5. Rootfs / mounts
- [ ] `./rootfs` path is hardcoded (`internal/container/rootfs.go:50`) — no way to
      point at a different image/rootfs per run
- [ ] Command to run is hardcoded to `/bin/ash` (`internal/container/child.go:15`) —
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
- [x] ~~Container identity: name/ID generation~~ — `internal/id` generates a
      random 256-bit hex ID per run, used as the container's hostname
- [ ] `vessel ps` / `vessel stop <id>` — now has metadata (`config.json`) and a
      real systemd unit name (`vessel-<id>.scope`) to build on; `ps` could
      list `/var/lib/vessel/containers/*`, `stop` could `systemctl stop` the
      matching scope directly instead of signaling a raw PID
- [ ] `vessel exec <id> <cmd>` — join an existing container's namespaces via `setns`
