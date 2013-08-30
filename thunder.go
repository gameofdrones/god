package main


import (
	"fmt"
	"log"

	"github.com/baloo/gousb/usb"
	"sync"
	"time"
)


const (
	X_TIME = 8000
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

type direction byte


const (
	DIR_UP = 0x00
	DIR_DOWN = 0x01
	DIR_LEFT = 0x02
	DIR_RIGHT = 0x03

	FIRE = 0x04 // TODO
	RELOAD = 0x05 // TODO

	WAIT = 0xFE
	STOP = 0xFF
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



func (c *Thunder) Move(dir direction) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	switch {
	case dir == DIR_UP:
		return c.Control(0x02)
	case dir == DIR_DOWN:
		return c.Control(0x01)
	case dir == DIR_LEFT:
		return nil
	case dir == DIR_RIGHT:
		return nil
	}

	return fmt.Errorf("thunder: you may specify only one direction") //You're insecure, Don't know what for, [...]
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
	ctype CommandType
	duration time.Duration
}

type CommandType byte



func (thunder *Thunder) RegisterMove(dir direction) error {
	log.Printf("register dir %s", dir)
	thunder.channel <- Command { ctype: (CommandType)( dir), duration: time.Since(time.Now()), }
	return nil
}

func (thunder *Thunder) RegisterStop() error {
	thunder.channel <- Command { ctype: STOP, duration: time.Since(time.Now()), }
	return nil
}

func (thunder *Thunder) RegisterFire() error {
	thunder.channel <- Command { ctype: FIRE, duration: time.Since(time.Now()), }
	return nil
}

func (thunder *Thunder) RegisterReload() error {
	thunder.channel <- Command { ctype: RELOAD, duration: time.Since(time.Now()), }
	return nil
}

func (thunder *Thunder) RegisterWait(duration time.Duration) error {
	thunder.channel <- Command { ctype: WAIT, duration: duration, }
	return nil
}

// to be launched as a goroutine
func (thunder *Thunder) Run() error {
	log.Printf("Run()")
	command := <- thunder.channel
	log.Printf("got command %s", command)
	switch {
	case command.ctype == FIRE:
		thunder.current_action = thunder.current_action | MOTOR_FIRE
		thunder.Control(thunder.current_action)
	case command.ctype == RELOAD:
		thunder.current_action = thunder.current_action & (0xFF -  MOTOR_FIRE)
		thunder.Control(thunder.current_action)
	case command.ctype == STOP:
		thunder.current_action = MOTOR_STOP
		thunder.Control(thunder.current_action)
		thunder.current_action = 0x00
		thunder.Control(thunder.current_action)
	case command.ctype == DIR_UP:
		thunder.current_action = (thunder.current_action & (0xFF - MOTOR_DOWN)) | MOTOR_UP
		thunder.Control(thunder.current_action)
	case command.ctype == DIR_DOWN:
		thunder.current_action = (thunder.current_action & (0xFF - MOTOR_UP)) | MOTOR_DOWN
		thunder.Control(thunder.current_action)
	case command.ctype == DIR_LEFT:
		thunder.current_action = (thunder.current_action & (0xFF - MOTOR_RIGHT)) | MOTOR_LEFT
		thunder.Control(thunder.current_action)
	case command.ctype == DIR_RIGHT:
		thunder.current_action = (thunder.current_action & (0xFF - MOTOR_LEFT)) | MOTOR_RIGHT
		thunder.Control(thunder.current_action)
	case command.ctype == WAIT:
		log.Printf("wait for %s", command.duration)
		time.Sleep(command.duration)
		log.Printf("awake")
	}
	return thunder.Run()
}


func (c *Thunder) SetPosition(x int, y int, fire bool) {
	c.mutex.Lock()

	// Check precision
	if c.moves >= 5 {
		c.dirty = true
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
	c.mutex.Unlock()
}


