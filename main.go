package main

import (
	"wifer/server/routes"
)

// Путь для os package никогда не должен начинаться с /
// Единственная причина, почему работает props.Conf.PATH - потому что данная переменная тоже не начинается с / :)

func main() {
	routes.Declare(&props)
	run()
}
