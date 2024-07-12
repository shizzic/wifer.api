package routes

import (
	"wifer/server/routes/other"
	"wifer/server/routes/user"
	"wifer/server/structs"
)

type Props = structs.Props

func Declare(props Props) {
	user.Declare(props)
	other.Declare(props)
}
