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
		ctx.RegisterMove(DIR_DOWN)
		ctx.RegisterWait(time.Duration(1)*time.Second)
		ctx.RegisterMove(DIR_UP)
		ctx.RegisterWait(time.Duration(1)*time.Second)
		ctx.RegisterMove(DIR_LEFT)
		ctx.RegisterWait(time.Duration(1)*time.Second)
		ctx.RegisterMove(DIR_RIGHT)
		ctx.RegisterWait(time.Duration(1)*time.Second)
		ctx.RegisterStop()
	}()
	go ctx.Run()
//	ctx.Move(DIR_DOWN)
//	ctx.Move(DIR_UP)
	

	_ = err

	time.Sleep(9 * 1e9)
}



