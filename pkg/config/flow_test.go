package config

import (
    "encoding/json"
    "testing"
)

func TestFlowJsonContainsGamepadAndTransmitter(t *testing.T) {
    flow := `{
  "nodes": [
    {"type":"channel","id":"30","data":{"label":"channel30"}},
    {"type":"channel","id":"31","data":{"label":"channel31"}},
    {"type":"channel","id":"32","data":{"label":"channel32"}},
    {"type":"channel","id":"33","data":{"label":"channel33"}},
    {"type":"channel","id":"34","data":{"label":"channel34"}},
    {"type":"tx","id":"35","data":{"label":"tx35","port":"/dev/ttyACM0"}},
    {"type":"axis","id":"9","data":{"label":"axis9","number":"0"}},
    {"type":"axis","id":"14","data":{"label":"axis9","number":"1"}},
    {"type":"axis","id":"19","data":{"label":"axis9","number":"2"}},
    {"type":"axis","id":"20","data":{"label":"axis9","number":"3"}},
    {"type":"axis","id":"21","data":{"label":"axis9","number":"4"}},
    {"type":"gamepad","id":"18","data":{"label":"gamepad18","name":"test","id":"r1","type":"udp"}}
  ],
  "edges": []
}`

    var doc struct {
        Nodes []struct {
            Type string          `json:"type"`
            Data json.RawMessage `json:"data"`
        } `json:"nodes"`
    }

    if err := json.Unmarshal([]byte(flow), &doc); err != nil {
        t.Fatalf("failed to parse flow json: %v", err)
    }

    foundTx := false
    foundGamepad := false
    for _, n := range doc.Nodes {
        if n.Type == "tx" {
            foundTx = true
        }
        if n.Type == "gamepad" {
            foundGamepad = true
        }
    }

    if !foundTx {
        t.Fatal("expected a transmitter node in the flow json")
    }
    if !foundGamepad {
        t.Fatal("expected a gamepad node in the flow json")
    }
}
