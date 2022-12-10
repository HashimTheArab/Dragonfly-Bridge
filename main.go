package main

import (
	"VBridge/commands"
	"VBridge/dfutils"
)

func main() {
	commands.RegisterAll()
	dfutils.StartServer()
}
