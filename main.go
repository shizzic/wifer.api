package main

import (
	"wifer/server/routes"
	"wifer/server/structs"
)

var props = Props{
	Conf: conf,
	Ctx:  ctx,
	DB:   DB,
	R:    router,
}

type Props = structs.Props
type Config = structs.Config
type Email = structs.Email

func main() {
	routes.Declare(&props)
	run()
}
