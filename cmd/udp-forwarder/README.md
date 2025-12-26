UDP Gamepad Forwarder

Build:

```bash
# linux/amd64
GOOS=linux GOARCH=amd64 go build -o udp-forwarder-linux ./cmd/udp-forwarder
# windows/amd64
GOOS=windows GOARCH=amd64 go build -o udp-forwarder.exe ./cmd/udp-forwarder
```

Usage:

- List controllers:

```bash
./udp-forwarder-linux --list
```

- Forward controller index 0 to UDP server at 192.168.1.10:9000:

```bash
./udp-forwarder-linux --index 0 --addr 192.168.1.10:9000
```

Notes:
- Requires SDL2 runtime on target machine (system package). On Linux install `libsdl2-dev`/`libsdl2-2.0` depending on distro. On Windows include SDL2 DLL next to the executable.
- The tool polls the joystick and sends JSON payloads similar to the main project's UDP format: `{"axes":{...},"buttons":{...},"hats":{...}}`.
