package main

import (
	"github.com/Youdontknowme720/Cimon/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
