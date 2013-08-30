package main


import (
	"fmt"
	"log"

	"github.com/baloo/gousb/usb"
	"sync"
	"time"
)


type direction byte

const (
	DIR_UP = 0x00
	DIR_DOWN = 0x01
	DIR_LEFT = 0x02
	DIR_RIGHT = 0x03

	FIRE = 0x04 // TODO

	WAIT = 0xFE
	STOP = 0xFF
)


type Thunder struct {
	usb_context *usb.Context
	device_id string
	subdevice *usb.Device
	mutex sync.Mutex
	channel chan Command
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
		thunder.Control(0x10)
	case command.ctype == STOP:
		thunder.Control(0x20)
	case command.ctype == DIR_UP:
		log.Printf("run dir up")
		thunder.Control(0x02)
	case command.ctype == DIR_DOWN:
		log.Printf("run dir down")
		thunder.Control(0x01)
	case command.ctype == DIR_LEFT:
		thunder.Control(0x04)
	case command.ctype == DIR_RIGHT:
		thunder.Control(0x08)
	case command.ctype == WAIT:
		log.Printf("wait for %s", command.duration)
		time.Sleep(command.duration)
		log.Printf("awake")
	}
	return thunder.Run()
}
