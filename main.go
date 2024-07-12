package main

import (
	"wifer/server/routes"
	"wifer/server/structs"
)

var props = Props{
	Ctx:  ctx,
	Conf: *conf,
	DB:   DB,
	R:    r,
}

type Props = structs.Props
type Config = structs.Config
type Email = structs.Email

func main() {
	routes.Declare(props)
	run()
}
