package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "net"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/veandco/go-sdl2/sdl"
)

type payload struct {
    Axes    map[string]int `json:"axes"`
    Buttons map[string]int `json:"buttons"`
    Hats    map[string]int `json:"hats"`
}

func listJoysticks() error {
    if err := sdl.Init(sdl.INIT_JOYSTICK); err != nil {
        return err
    }
    defer sdl.Quit()

    n := sdl.NumJoysticks()
    fmt.Printf("Found %d joystick(s)\n", n)
    for i := 0; i < n; i++ {
        name := sdl.JoystickNameForIndex(i)
        fmt.Printf("%d: %s\n", i, name)
    }
    return nil
}

func forward(index int, addr string, interval time.Duration) error {
    if err := sdl.Init(sdl.INIT_JOYSTICK); err != nil {
        return err
    }
    defer sdl.Quit()

    n := sdl.NumJoysticks()
    if n == 0 {
        return fmt.Errorf("no joysticks found")
    }
    if index < 0 || index >= n {
        return fmt.Errorf("joystick index out of range: %d (found %d)", index, n)
    }

    j := sdl.JoystickOpen(index)
    if j == nil {
        return fmt.Errorf("could not open joystick %d", index)
    }
    defer j.Close()

    udpAddr, err := net.ResolveUDPAddr("udp", addr)
    if err != nil {
        return err
    }
    conn, err := net.DialUDP("udp", nil, udpAddr)
    if err != nil {
        return err
    }
    defer conn.Close()

    sig := make(chan os.Signal, 1)
    signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

    ticker := time.NewTicker(interval)
    defer ticker.Stop()


    for {
        select {
        case <-sig:
            fmt.Println("shutting down")
            return nil
        case <-ticker.C:
            // pump SDL events so joystick state is up-to-date
            sdl.PumpEvents()

            // poll joystick
            axes := make(map[string]int)
            buttons := make(map[string]int)
            hats := make(map[string]int)

            for a := 0; a < int(j.NumAxes()); a++ {
                v := int(j.Axis(a))
                axes[fmt.Sprintf("%d", a)] = v
            }
            for b := 0; b < int(j.NumButtons()); b++ {
                vb := int(j.Button(b))
                buttons[fmt.Sprintf("%d", b)] = vb
            }
            for h := 0; h < int(j.NumHats()); h++ {
                hv := int(j.Hat(h))
                // normalize hat values to -1/0/1 for simplicity
                val := 0
                if hv == sdl.HAT_UP || hv == sdl.HAT_RIGHTUP || hv == sdl.HAT_LEFTUP {
                    val = 1
                } else if hv == sdl.HAT_DOWN || hv == sdl.HAT_RIGHTDOWN || hv == sdl.HAT_LEFTDOWN {
                    val = -1
                }
                hats[fmt.Sprintf("%d", h)] = val
            }

            p := payload{Axes: axes, Buttons: buttons, Hats: hats}
            data, _ := json.Marshal(p)
            _, _ = conn.Write(data)
        }
    }
}

func main() {
    addr := flag.String("addr", "127.0.0.1:9000", "UDP server address to send inputs to")
    list := flag.Bool("list", false, "List available joysticks and exit")
    index := flag.Int("index", 0, "Joystick index to forward (from --list output)")
    interval := flag.Duration("interval", 20*time.Millisecond, "Polling interval")
    hz := flag.Int("hz", 0, "Polling rate in Hz (if >0 overrides --interval)")
    flag.Parse()

    if *list {
        if err := listJoysticks(); err != nil {
            fmt.Fprintf(os.Stderr, "error: %v\n", err)
            os.Exit(1)
        }
        return
    }

    // compute effective interval
    effInterval := *interval
    if *hz > 0 {
        effInterval = time.Second / time.Duration(*hz)
    }

    fmt.Printf("forwarding joystick %d -> %s (interval %s)\n", *index, *addr, effInterval.String())
    if err := forward(*index, *addr, effInterval); err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }
}
