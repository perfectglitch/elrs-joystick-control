package devices

import (
    "encoding/json"
    "fmt"
    "net"
    "sync"
    "time"
    "github.com/kaack/elrs-joystick-control/pkg/util"
)

// UDPGamepad receives gamepad-like input over UDP (JSON) and implements Gamepad.
type UDPGamepad struct {
    Id   string `json:"id"`
    Name string `json:"name"`

    mu      sync.RWMutex
    axes    map[int]util.RawValue
    buttons map[int]util.RawValue
    hats    map[int]util.RawValue

    axesCount    int32
    buttonsCount int32
    hatsCount    int32

    conn *net.UDPConn
    quit chan struct{}
    last time.Time
}

// NewUDPGamepad creates and starts a UDP listener on the given address (e.g. ":9000").
func NewUDPGamepad(id, name, addr string) (*UDPGamepad, error) {
    udp := &UDPGamepad{
        Id:   id,
        Name: name,
        axes: make(map[int]util.RawValue),
        buttons: make(map[int]util.RawValue),
        hats: make(map[int]util.RawValue),
        quit: make(chan struct{}),
    }

    udpAddr, err := net.ResolveUDPAddr("udp", addr)
    if err != nil {
		fmt.Printf("Error resolving UDP address %s: %v", addr, err)
        return nil, err
    }

    conn, err := net.ListenUDP("udp", udpAddr)
    if err != nil {
		fmt.Printf("Error listening on UDP %s: %v\n", addr, err)
        return nil, err
    }
    udp.conn = conn

    go udp.readLoop()

    return udp, nil
}

// Close implements Gamepad.Close
func (u *UDPGamepad) Close() {
    close(u.quit)
    if u.conn != nil {
        _ = u.conn.Close()
    }
}

func (u *UDPGamepad) InstanceId() int32 { return 0 }
func (u *UDPGamepad) Axes() int32 {
    return u.axesCount
}
func (u *UDPGamepad) Buttons() int32 { return u.buttonsCount }
func (u *UDPGamepad) Hats() int32 { return u.hatsCount }

func (u *UDPGamepad) Axis(axis int) util.RawValue {
    u.mu.RLock()
    defer u.mu.RUnlock()
    if v, ok := u.axes[axis]; ok {
        return v
    }
    return 0
}

func (u *UDPGamepad) Button(button int) util.RawValue {
    u.mu.RLock()
    defer u.mu.RUnlock()
    if v, ok := u.buttons[button]; ok {
        return v
    }
    return 0
}

func (u *UDPGamepad) Hat(hat int) util.RawValue {
    u.mu.RLock()
    defer u.mu.RUnlock()
    if v, ok := u.hats[hat]; ok {
        return util.MapRange(v, -1, 1, util.MinRaw, util.MaxRaw)
    }
    return 0
}

type udpPayload struct {
    Axes    map[string]util.RawValue `json:"axes"`
    Buttons map[string]util.RawValue `json:"buttons"`
    Hats    map[string]util.RawValue `json:"hats"`
}

func (u *UDPGamepad) readLoop() {
    buf := make([]byte, 2048)
    for {
        select {
        case <-u.quit:
            return
        default:
            u.conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
            n, remoteAddr, err := u.conn.ReadFromUDP(buf)
            if err != nil {
                if ne, ok := err.(net.Error); ok && ne.Timeout() {
                    continue
                }
                fmt.Printf("UDP read error: %v\n", err)
                continue
            }

            payload := buf[:n]
            var p udpPayload
            if err := json.Unmarshal(payload, &p); err != nil {
                fmt.Printf("UDP[%s] parse error: %v\n", remoteAddr.String(), err)
                continue
            }

            u.mu.Lock()
            for k, v := range p.Axes {
                // keys are strings with integers
                var idx int
                _, _ = fmt.Sscanf(k, "%d", &idx)
                u.axes[idx] = v
                if idx+1 > int(u.axesCount) {
                    u.axesCount = int32(idx + 1)
                }
            }
            for k, v := range p.Buttons {
                var idx int
                _, _ = fmt.Sscanf(k, "%d", &idx)
                u.buttons[idx] = v
                if idx+1 > int(u.buttonsCount) {
                    u.buttonsCount = int32(idx + 1)
                }
            }
            for k, v := range p.Hats {
                var idx int
                _, _ = fmt.Sscanf(k, "%d", &idx)
                u.hats[idx] = v
                if idx+1 > int(u.hatsCount) {
                    u.hatsCount = int32(idx + 1)
                }
            }
            u.last = time.Now()
            u.mu.Unlock()

            // notify interested parties (e.g. config eval loop) that device data changed
            notifyDeviceEvent()
        }
    }
}

// MarshalJSON makes UDPGamepad marshal similarly to InputGamepad for API compatibility.
func (u *UDPGamepad) MarshalJSON() ([]byte, error) {
    type fake struct {
        Id   string `json:"id"`
        Name string `json:"name"`
    }
    return json.Marshal(struct {
        fake
        Axes    int32 `json:"axes"`
        Buttons int32 `json:"buttons"`
        Hats    int32 `json:"hats"`
    }{
        fake:    fake{Id: u.Id, Name: u.Name},
        Axes:    u.Axes(),
        Buttons: u.Buttons(),
        Hats:    u.Hats(),
    })
}
