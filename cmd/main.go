package main

import (
	"flag"
	"fmt"
	"wlk_tlk/internal/stantion"
)

func main() {

	fmt.Println("Hello world to webrtc")
	port := flag.Int("port",8070,"port of ws station")

	flag.Parse()

	stantion := stantion.NewStantion()
	stantion.Run(*port)
}