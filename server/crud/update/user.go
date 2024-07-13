package update

import (
	"net/http"
	"time"
	"wifer/server/auth"
	"wifer/server/structs"

	"go.mongodb.org/mongo-driver/bson"
)

type Props = structs.Props

// Изменить дату последнего посещения для пользователя
func ChangeLastOnline(props *Props, timestamp bool, id int) {
	props.DB["users"].UpdateOne(props.Ctx, bson.M{"_id": id}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "online", Value: timestamp}}},
		{Key: "$set", Value: bson.D{{Key: "last_time", Value: time.Now().Unix()}}},
	})
}

func Logout(w http.ResponseWriter, r *http.Request, props *Props, id int) {
	props.DB["users"].UpdateOne(props.Ctx, bson.M{"_id": id}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "online", Value: false}}},
		{Key: "$set", Value: bson.D{{Key: "last_time", Value: time.Now().Unix()}}},
	})

	auth.MakeCookies(props, "", "", -1, w)
	http.SetCookie(w, &http.Cookie{
		Name:     "premium",
		Value:    "premium",
		Path:     "/",
		Domain:   "." + props.Conf.SELF_DOMAIN_NAME,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
}
