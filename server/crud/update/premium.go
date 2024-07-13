package update

import (
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Premium(props *Props, w http.ResponseWriter, r *http.Request, id int, user primitive.M) {
	// Если у пользователя был активирован премиум когда либо до и все еще не обнулен
	if user["premium"].(int64) != 0 {
		// Если премиум просрочился, обнуляю в бд и удаляю куку
		if user["premium"].(int64) <= time.Now().Unix() {
			props.DB["users"].UpdateOne(props.Ctx, bson.M{"_id": id}, bson.D{
				{Key: "$set", Value: bson.D{{Key: "premium", Value: 0}}},
			})

			user["premium"] = 0

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
		} else {
			// Создаю куку, если ее не было до
			if _, err := r.Cookie("premium"); err != nil {
				http.SetCookie(w, &http.Cookie{
					Name:     "premium",
					Value:    "premium",
					Path:     "/",
					Domain:   "." + props.Conf.SELF_DOMAIN_NAME,
					MaxAge:   int(user["premium"].(int64) - time.Now().Unix()),
					HttpOnly: true,
					Secure:   true,
					SameSite: http.SameSiteNoneMode,
				})
			}
		}
	}
}
