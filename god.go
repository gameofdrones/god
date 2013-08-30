package main

import (
	"flag"
	"log"

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



	_ = err

}



