package main


import (
	"fmt"

	"github.com/baloo/gousb/usb"
)


type direction byte

const (
	DIR_UP = 0x00
	DIR_DOWN = 0x01
	DIR_LEFT = 0x02
	DIR_RIGHT = 0x03
)


type Thunder struct {
	usb_context *usb.Context
	device_id string
	subdevice *usb.Device
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
	switch {
	case dir == DIR_UP:
		return nil
	case dir == DIR_DOWN:
		return nil
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



