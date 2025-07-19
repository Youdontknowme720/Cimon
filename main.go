package main

import "gitlab.com/ayan0k0uji-group/Cimon/cmd"

func main() {
    if err := cmd.Execute(); err != nil {
        panic(err)
    }
}