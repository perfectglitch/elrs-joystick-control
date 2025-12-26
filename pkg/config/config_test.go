package config

import (
    "testing"
)

func TestConfigContainsGamepadAndTransmitter(t *testing.T) {
    js := `{
        "config": {
            "input_output_map": {
                "tx1": {
                    "id": "tx1",
                    "type": "tx",
                    "tx": {
                        "port": "/dev/ttyACM0",
                        "channels": [{"id":"ch1","type":"channel","channel":{"number":1,"input":{"id":"tx1","type":"tx"}}}]
                    }
                },
                "gp1": {
                    "id": "gp1",
                    "type": "gamepad",
                    "gamepad": {
                        "id": "r1",
                        "name": "test",
                        "type": "udp",
                        "udp_addr": ":9000"
                    }
                }
            }
        }
    }`

    var ctl Controller
    if err := ctl.UnmarshalJSON([]byte(js)); err != nil {
        t.Fatalf("failed to unmarshal config: %v", err)
    }

    if ctl.Config == nil || ctl.Config.IOMap == nil {
        t.Fatal("config or IOMap is nil after unmarshal")
    }

    hasGamepad := false
    hasTx := false
    for k, ih := range ctl.Config.IOMap {
        if ih == nil || ih.IO == nil {
            continue
        }
        switch ih.IO.(type) {
        case *InputGamepad:
            hasGamepad = true
        case *OutputTransmitter:
            hasTx = true
        }
        _ = k
    }

    if !hasTx {
        t.Fatal("expected transmitter in IOMap but not found")
    }
    if !hasGamepad {
        t.Fatal("expected gamepad in IOMap but not found")
    }
}
