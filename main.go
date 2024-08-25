package main

import (
	"os"
)

func main() {
	os.Create("test.txt")
	// dump.PrepareDB(&props)
	// routes.Declare(&props)
	run()
}
