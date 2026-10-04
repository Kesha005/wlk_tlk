package main

import (
	"fmt"
	"wlk_tlk/internal/stantion"
)

func main() {

	fmt.Println("Hello world to webrtc")
	stantion := stantion.NewStantion()
	stantion.Run(8070)
}