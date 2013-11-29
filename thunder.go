package main


import (
        "fmt"
        "log"
        "os/exec"

        "github.com/baloo/gousb/usb"
        "sync"
        "time"
)


const (
        X_TIME = 5360
        X_POSITIONS = 100
        Y_TIME = 1000
        Y_POSITIONS = 10
)

const (
        MOTOR_UP = 0x02
        MOTOR_DOWN = 0x01
        MOTOR_LEFT = 0x04
        MOTOR_RIGHT = 0x08
        MOTOR_FIRE = 0x10
        MOTOR_STOP = 0x20
)

type Action byte
const (
        DIR_UP Action = 0x00
        DIR_DOWN Action = 0x01
        DIR_LEFT Action = 0x02
        DIR_RIGHT Action = 0x03

        FIRE Action = 0x04 // TODO
        RELOAD Action = 0x05 // TODO

        WAIT Action = 0xFE
        STOP Action = 0xFF
)


type Thunder struct {
        usb_context *usb.Context
        device_id string
        subdevice *usb.Device
        mutex sync.Mutex
        channel chan Command
        current_action byte
        current_x int
        current_y int
        dirty bool
        moves int
}



func NewThunder(device_id string) (*Thunder, error) {
        usb_context := usb.NewContext()
        dev, err := FindDevice(usb_context, device_id)

        if err != nil {
                return nil, err
        }


        c := &Thunder{
                usb_context: usb_context,
                device_id: device_id,
                subdevice: dev,
                channel: make(chan Command),
                dirty: true,
        }

        return c, nil
}

func (c *Thunder) Close() error {
        if c.subdevice != nil {
                c.subdevice.Close()
        }
        c.usb_context.Close()

        return nil
}

func FindDevice(ctx *usb.Context, device string) (*usb.Device, error) {
        // Lookup for devices
        devs, err := ctx.ListDevices(func(desc *usb.Descriptor) bool {
                if fmt.Sprintf("%s:%s", desc.Vendor, desc.Product) != device {
                        return false
                }

                return true
        })

        if err != nil {
                return nil, fmt.Errorf("thunder: Error while looking for devices. %s", err)
        }

        if len(devs) == 0 {
                return nil, fmt.Errorf("thunder: Unable to find any devices")
        }

        // Keep one device
        dev := devs[0]
        devs = devs[1:]

        // Close any unused devices
        defer func() {
                for _, d := range devs {
                        d.Close()
                }
        }()

        dev.DetachKernelDriver(0)

        return dev, nil
}

func (c *Thunder) Stop() error {
        return c.Control(0x20)
}

func (c *Thunder) Control(msg byte) error {
        log.Printf("Current state: %d", c.current_action)
        data := []byte{0x02,msg,0x00,0x00,0x00,0x00,0x00,0x00}

        ep, err :=  c.subdevice.Control(
                0x21,
                0x09, //request
                0x00, //wvalue
                0x00, //windex
                data)

        _ = ep
        _ = err

        return nil // TODO
}


type Command struct {
        action Action
        duration time.Duration
}

type CommandType byte



func (thunder *Thunder) RegisterMove(action Action) error {
        log.Printf("register action %s", action)
        thunder.channel <- Command { action: action, duration: time.Since(time.Now()), }
        return nil
}

func (thunder *Thunder) RegisterStop() error {
        thunder.channel <- Command { action: STOP, duration: time.Since(time.Now()), }
        return nil
}

func (thunder *Thunder) RegisterFire() error {
        thunder.channel <- Command { action: FIRE, duration: time.Since(time.Now()), }
        return nil
}

func (thunder *Thunder) RegisterReload() error {
        thunder.channel <- Command { action: RELOAD, duration: time.Since(time.Now()), }
        return nil
}

func (thunder *Thunder) RegisterWait(duration time.Duration) error {
        thunder.channel <- Command { action: WAIT, duration: duration, }
        return nil
}

// to be launched as a goroutine
func (thunder *Thunder) Run() error {
        log.Printf("Run()")
        command := <- thunder.channel
        log.Printf("got command %s", command)
        switch(command.action) {
        case FIRE:
                thunder.current_action = thunder.current_action | MOTOR_FIRE
                thunder.Control(thunder.current_action)
                GrabMjpegFrame()
        case RELOAD:
                thunder.current_action = thunder.current_action & (0xFF -  MOTOR_FIRE)
                thunder.Control(thunder.current_action)
        case STOP:
                thunder.current_action = MOTOR_STOP
                thunder.Control(thunder.current_action)
                thunder.current_action = 0x00
                thunder.Control(thunder.current_action)
        case DIR_UP:
                thunder.current_action = (thunder.current_action & (0xFF - MOTOR_DOWN)) | MOTOR_UP
                thunder.Control(thunder.current_action)
        case DIR_DOWN:
                thunder.current_action = (thunder.current_action & (0xFF - MOTOR_UP)) | MOTOR_DOWN
                thunder.Control(thunder.current_action)
        case DIR_LEFT:
                thunder.current_action = (thunder.current_action & (0xFF - MOTOR_RIGHT)) | MOTOR_LEFT
                thunder.Control(thunder.current_action)
        case DIR_RIGHT:
                thunder.current_action = (thunder.current_action & (0xFF - MOTOR_LEFT)) | MOTOR_RIGHT
                thunder.Control(thunder.current_action)
        case WAIT:
                log.Printf("wait for %s", command.duration)
                time.Sleep(command.duration)
                log.Printf("awake")
        }
        return thunder.Run()
}


func (c *Thunder) ResetPosition() {
        c.mutex.Lock()
        defer c.mutex.Unlock()

        c.RegisterMove(DIR_DOWN)
        c.RegisterMove(DIR_LEFT)
        c.RegisterWait(time.Duration(8)*time.Second)
        c.RegisterStop()
        c.dirty = false
        c.current_x = 0
        c.current_y = 0
}
func (c *Thunder) SetPosition(x int, y int, fire bool) {
        c.mutex.Lock()
        defer c.mutex.Unlock()

        // Check precision
        if c.moves >= 5 {
                c.dirty = true
                c.moves = 0
        }

        // Check we did not move (c) John Drummond
        if c.dirty {
                c.RegisterMove(DIR_DOWN)
                c.RegisterMove(DIR_LEFT)
                c.RegisterWait(time.Duration(8)*time.Second)
                c.RegisterStop()
                c.dirty = false
                c.current_x = 0
                c.current_y = 0
        }

        delta_x := x - c.current_x
        c.current_x = x
        log.Printf("will delta x: %d", delta_x)
        if delta_x > 0 {
                wait := X_TIME / X_POSITIONS * delta_x
                c.RegisterMove(DIR_RIGHT)
                c.RegisterWait(time.Duration(wait)*time.Millisecond)
        }
        if delta_x < 0 {
                wait := X_TIME / X_POSITIONS * delta_x * -1
                c.RegisterMove(DIR_LEFT)
                c.RegisterWait(time.Duration(wait)*time.Millisecond)
        }
        c.RegisterStop()

        delta_y := y - c.current_y
        c.current_y = y
        log.Printf("will delta y: %d", delta_y)
        if delta_y > 0 {
                wait := Y_TIME / Y_POSITIONS * delta_y
                c.RegisterMove(DIR_UP)
                c.RegisterWait(time.Duration(wait)*time.Millisecond)
        }
        if delta_y < 0 {
                wait := Y_TIME / Y_POSITIONS * delta_y * -1
                c.RegisterMove(DIR_DOWN)
                c.RegisterWait(time.Duration(wait)*time.Millisecond)
        }

        c.RegisterStop()

        if fire {
                c.RegisterFire()
                c.RegisterWait(time.Duration(5000)*time.Millisecond)
                c.RegisterReload()
        }

        c.moves = c.moves + 1
}


func (thunder *Thunder) Put(action Action) error {
        thunder.mutex.Lock()
        defer thunder.mutex.Unlock()
        thunder.RegisterMove(action)
        thunder.dirty = true
        return nil
}

func (thunder *Thunder) Fire() error {
        thunder.mutex.Lock()
        defer thunder.mutex.Unlock()
        thunder.RegisterFire()
        thunder.RegisterWait(time.Duration(5000)*time.Millisecond)
        thunder.RegisterReload()
        return nil
}

func (thunder *Thunder) Delete(action Action) error {
        thunder.mutex.Lock()
        defer thunder.mutex.Unlock()
        thunder.current_action = (thunder.current_action & (0xFF - (byte)(action)))
        thunder.Control(thunder.current_action)
        thunder.dirty = true
        return nil
}

func GrabMjpegFrame() {
        shotImg := "lastShot.jpg"
        hostWithPort := "10.0.25.113:8081"
        scriptName := "grab_mjpeg_frame.py"
        cmd := exec.Command("python", scriptName, hostWithPort, shotImg)
        err := cmd.Start()
        if err != nil {
                log.Printf("2")
                log.Fatal(err)
        }
        err = cmd.Wait()
        if err != nil {
                log.Printf("Error while calling the script. Error: %v\n", err)
        } else {
                log.Printf("Image successfully created.")
        }

        thunder.current_action = (thunder.current_action & (0xFF - (byte)(action)))
        thunder.Control(thunder.current_action)
        return nil
}
