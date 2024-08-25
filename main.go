package main

import "wifer/server/routes"

func main() {
	routes.Declare(&props)
	run()
}
