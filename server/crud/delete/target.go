package delete

import (
	"wifer/server/structs"

	"go.mongodb.org/mongo-driver/bson"
)

func TargetLike(props *structs.Props, data *structs.Target, id int) {
	if id > 0 && id != data.Target && data.Target > 0 {
		props.DB["likes"].DeleteOne(props.Ctx, bson.M{"user": id, "target": data.Target})
	}
}

func TargetPrivate(props *structs.Props, data *structs.Target, id int) {
	if id > 0 && id != data.Target && data.Target > 0 {
		props.DB["private"].DeleteOne(props.Ctx, bson.M{"user": id, "target": data.Target})
	}
}

func TargetAccess(props *structs.Props, data *structs.Target, id int) {
	if id > 0 && id != data.Target && data.Target > 0 {
		props.DB["access"].DeleteOne(props.Ctx, bson.M{"user": id, "target": data.Target})
	}
}
