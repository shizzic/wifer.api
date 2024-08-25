package main

import (
	"wifer/cron/dump"
	"wifer/server/routes"
)

func main() {
	dump.PrepareDB(&props)
	routes.Declare(&props)
	run()
}
