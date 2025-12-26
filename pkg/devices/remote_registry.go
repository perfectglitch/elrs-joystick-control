// SPDX-FileCopyrightText: © 2023 OneEyeFPV oneeyefpv@gmail.com
// SPDX-License-Identifier: GPL-3.0-or-later
// SPDX-License-Identifier: FS-0.9-or-later

package devices

import "sync"

var (
    remoteGamepadsMu sync.RWMutex
    remoteGamepads   = map[string]Gamepad{}
    deviceEventCbMu  sync.RWMutex
    deviceEventCb    func()
)

// RegisterRemoteGamepad registers a remote gamepad in the global registry.
func RegisterRemoteGamepad(id string, g Gamepad) {
    remoteGamepadsMu.Lock()
    defer remoteGamepadsMu.Unlock()
    remoteGamepads[id] = g
}

// UnregisterRemoteGamepad removes a gamepad from the global registry.
func UnregisterRemoteGamepad(id string) {
    remoteGamepadsMu.Lock()
    defer remoteGamepadsMu.Unlock()
    delete(remoteGamepads, id)
}

// GetRemoteGamepads returns a copy of the registered remote gamepads.
func GetRemoteGamepads() map[string]Gamepad {
    remoteGamepadsMu.RLock()
    defer remoteGamepadsMu.RUnlock()
    res := make(map[string]Gamepad, len(remoteGamepads))
    for k, v := range remoteGamepads {
        res[k] = v
    }
    return res
}

// SetDeviceEventCallback registers a callback that will be invoked when a remote
// device (like UDPGamepad) receives new input. Pass nil to clear.
func SetDeviceEventCallback(cb func()) {
    deviceEventCbMu.Lock()
    defer deviceEventCbMu.Unlock()
    deviceEventCb = cb
}

// notifyDeviceEvent calls the registered callback if set.
func notifyDeviceEvent() {
    deviceEventCbMu.RLock()
    cb := deviceEventCb
    deviceEventCbMu.RUnlock()
    if cb != nil {
        cb()
    }
}
