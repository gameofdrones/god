package main

import (
	"flag"
	"log"

	"time"
)


var (
	debug = flag.Int("debug", 0, "libusb debug level (0..3)")
)

func main() {
	log.Printf("Starting god")

	device := "2123:1010"

	// Create a new context for managing thunder rocket launcher
	ctx, err := NewThunder(device)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	// Ensure everything is closed on exit
	defer ctx.Close()

	go func() {
		//ctx.RegisterFire()
		ctx.SetPosition(20, 5, true)
		ctx.SetPosition(30, 0, true)
		ctx.SetPosition(10, 10, true)
		ctx.SetPosition(100, 10, true)
		ctx.SetPosition(90, 10, true)
		ctx.SetPosition(30, 10, true)
		//ctx.RegisterStop()
		//ctx.RegisterMove(DIR_DOWN)
		//ctx.RegisterMove(DIR_UP)
		//ctx.RegisterWait(time.Duration(15000)*time.Millisecond)
		//ctx.RegisterMove(DIR_LEFT)
		//ctx.RegisterWait(time.Duration(500)*time.Millisecond)
		//ctx.RegisterMove(DIR_UP)
		//ctx.RegisterWait(time.Duration(500)*time.Millisecond)
		//ctx.RegisterMove(DIR_DOWN)
		//ctx.RegisterWait(time.Duration(500)*time.Millisecond)
		//ctx.RegisterMove(DIR_UP)
		//ctx.RegisterWait(time.Duration(500)*time.Millisecond)
		//ctx.RegisterMove(DIR_DOWN)
		//ctx.RegisterReload()
		//ctx.RegisterFire()
		//ctx.RegisterWait(time.Duration(1500)*time.Millisecond)
		//ctx.RegisterMove(DIR_UP)
		//ctx.RegisterWait(time.Duration(500)*time.Millisecond)
		//ctx.RegisterMove(DIR_DOWN)
		//ctx.RegisterWait(time.Duration(500)*time.Millisecond)
		//ctx.RegisterReload()
		//ctx.RegisterStop()
	}()
	go ctx.Run()
//	ctx.Move(DIR_DOWN)
//	ctx.Move(DIR_UP)
	

	_ = err

	time.Sleep(100 * 1e9)
}



