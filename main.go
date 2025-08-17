package main

import (
	"fmt"

	"github.com/Youdontknowme720/Cimonv2/ui"
)

func main() {
	fmt.Print("Start App")
	myApp := ui.NewApp()
	myApp.Setup()
	err := myApp.Run()
	if err != nil {
		return
	}
}
