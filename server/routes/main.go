package routes

// Вывести мапу с ключами
// fmt.Printf("%+v", data)

import (
	"wifer/server/structs"

	unrolled "github.com/unrolled/render"
)

type Props = structs.Props
type Auth = structs.Auth
type User = structs.User
type Signin = structs.Signin
type Images = structs.Images

var render = unrolled.New()

func Declare(props *Props) {
	user(props)
	location(props)
	image(props)
	other(props)
}
