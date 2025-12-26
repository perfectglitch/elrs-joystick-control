package devices

import (
    "encoding/json"
    "net"
    "testing"
    "time"

    "github.com/kaack/elrs-joystick-control/pkg/util"
)

func TestUDPGamepadReceivesAndMarshals(t *testing.T) {
    u, err := NewUDPGamepad("testid", "testname", ":0")
    if err != nil {
        t.Fatalf("NewUDPGamepad error: %v", err)
    }
    defer u.Close()

    // send a test packet to the listening address
    addr := u.conn.LocalAddr().(*net.UDPAddr)
    c, err := net.DialUDP("udp", nil, addr)
    if err != nil {
        t.Fatalf("DialUDP error: %v", err)
    }
    defer c.Close()

    payload := map[string]map[string]util.RawValue{
        "axes": {
            "0": util.RawValue(12345),
        },
        "buttons": {
            "1": util.RawValue(1),
        },
        "hats": {
            "0": util.RawValue(0),
        },
    }

    b, _ := json.Marshal(payload)
    if _, err := c.Write(b); err != nil {
        t.Fatalf("write error: %v", err)
    }

    // give the reader goroutine some time
    time.Sleep(150 * time.Millisecond)

    if got := u.Axis(0); got != util.RawValue(12345) {
        t.Fatalf("Axis(0) = %v, want %v", got, util.RawValue(12345))
    }

    if got := u.Button(1); got != util.RawValue(1) {
        t.Fatalf("Button(1) = %v, want %v", got, util.RawValue(1))
    }

    if got := u.Hat(0); got != util.RawValue(0) {
        t.Fatalf("Hat(0) = %v, want %v", got, util.RawValue(0))
    }

    // test JSON marshaling reports counts
    jb, err := json.Marshal(u)
    if err != nil {
        t.Fatalf("marshal error: %v", err)
    }
    var jm map[string]any
    if err := json.Unmarshal(jb, &jm); err != nil {
        t.Fatalf("unmarshal marshaled json error: %v", err)
    }

    if axes, ok := jm["axes"].(float64); !ok || int32(axes) < 1 {
        t.Fatalf("unexpected axes in marshaled json: %v", jm["axes"])
    }
    if buttons, ok := jm["buttons"].(float64); !ok || int32(buttons) < 1 {
        t.Fatalf("unexpected buttons in marshaled json: %v", jm["buttons"])
    }
}
