package update

import (
	"errors"
	"net/http"
	"time"
	"wifer/server/structs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Актуализирую куку премиума. Возвращает наличие премиума
func Premium(props *Props, w http.ResponseWriter, r *http.Request, id int, user primitive.M) bool {
	// Если у пользователя был активирован премиум когда либо до и все еще не обнулен
	if user["premium"].(int64) != 0 {
		// Если премиум просрочился, обнуляю в бд и удаляю куку
		if user["premium"].(int64) <= time.Now().Unix() {
			props.DB["users"].UpdateOne(props.Ctx, bson.M{"_id": id}, bson.D{
				{Key: "$set", Value: bson.D{{Key: "premium", Value: int64(0)}}},
			})

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

			return false
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

			return true
		}
	}

	return false
}

func ActivateOneTimeTrial(props *structs.Props, w http.ResponseWriter, id int) (int64, error) {
	var user bson.M
	opts := options.FindOne().SetProjection(bson.M{"_id": 0, "trial": 1, "premium": 1})
	props.DB["users"].FindOne(props.Ctx, bson.M{"_id": id}, opts).Decode(&user)

	// Даю пробник, только если он не был активирован до этого
	if !user["trial"].(bool) {
		now := time.Now().Unix()
		plus := int64(86400 * 7) // кол-во дней, которое я добавляю
		expires := now + plus    // добавляю 7 дней от нынешнего момента
		var maxAge int           // это для куки

		if user["premium"].(int64) == 0 {
			maxAge = int(plus)
			props.DB["users"].UpdateOne(props.Ctx, bson.M{"_id": id}, bson.D{
				{Key: "$set", Value: bson.D{{Key: "trial", Value: true}}},
				{Key: "$set", Value: bson.D{{Key: "premium", Value: expires}}},
			})
		} else {
			expires = int64(user["premium"].(int64) + plus) // добавляю уже имеющийся примиум (его остаток)
			maxAge = int(user["premium"].(int64) - now + plus)
			props.DB["users"].UpdateOne(props.Ctx, bson.M{"_id": id}, bson.D{
				{Key: "$set", Value: bson.D{{Key: "trial", Value: true}}},
				{Key: "$set", Value: bson.D{{Key: "premium", Value: expires}}}, // прибавляю 7 дней
			})
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "premium",
			Value:    "premium",
			Path:     "/",
			Domain:   "." + props.Conf.SELF_DOMAIN_NAME,
			MaxAge:   maxAge,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		})
		return expires, nil
	}

	return 0, errors.New("activated_before")
}
