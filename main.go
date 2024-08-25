package main

import (
	"wifer/server/routes"
)

// Путь для os package никогда не должен начинаться с /
// Единственная причина, почему работает props.Conf.PATH - потому что данная переменная тоже не начинается с / :)
// Я все равно добавляю props.Conf.PATH везде, для кросс платформенности

func main() {
	routes.Declare(&props)
	run()
}
