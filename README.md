# Vessel

Vessel is a lightweight container written in Go for learning how Linux containers work under the hood.

The project aims to build a minimal container from scratch without relying on existing container engines. It focuses on understanding Linux primitives such as:

- Linux namespaces
- Root filesystems
- Cgroups

> **Note:** This is an educational project and is **not** intended for production use.

## Prerequisites

- Linux
- Go 1.26+ (or your installed Go version)

## Getting Started

### 1. Download the Alpine miniroot filesystem

```bash
./scripts/setup-rootfs.sh
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

## License

This project is provided for educational purposes.