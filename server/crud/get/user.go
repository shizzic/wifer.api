package get

import (
	"errors"
	"net/http"
	"strconv"
	"wifer/server/auth"
	"wifer/server/crud/update"
	"wifer/server/structs"

	unrolled "github.com/unrolled/render"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Signin = structs.Signin

var render = unrolled.New()

// Получаем id пользователя из куки, для авторизации
func UserID(w http.ResponseWriter, r *http.Request, props *Props) (result int) {
	id, err := r.Cookie("id")
	if err != nil {
		render.JSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	result, _ = strconv.Atoi(id.Value)
	return
}

// Получаем весь профиль пользователя
func Profile(id int, props *Props) (bson.M, error) {
	var user bson.M
	opts := options.FindOne().SetProjection(bson.M{
		"username":   1,
		"title":      1,
		"about":      1,
		"sex":        1,
		"age":        1,
		"body":       1,
		"height":     1,
		"weight":     1,
		"smokes":     1,
		"drinks":     1,
		"ethnicity":  1,
		"search":     1,
		"income":     1,
		"children":   1,
		"industry":   1,
		"premium":    1,
		"avatar":     1,
		"public":     1,
		"private":    1,
		"prefer":     1,
		"created_at": 1,
		"last_time":  1,
		"online":     1,
		"country_id": 1,
		"city_id":    1,
	})

	if err := props.DB["users"].FindOne(props.Ctx, bson.M{"_id": id, "status": true}, opts).Decode(&user); err != nil {
		return user, errors.New("0")
	}

	return user, nil
}

func UserEmailByApi(data Signin) (email string, err error) {
	switch data.Method {
	case "Google":
		email, err = auth.IsGoogle(data.ID, data.Token)
	case "Facebook":
		email, err = auth.IsFacebook(data.ID, data.Token)
	}
	return
}

// Получаю пользователя и список айдишников юзеров, которые ему отписали (не прочитанные сообщения)
func UserAndMessagedHimIds(props *Props, w http.ResponseWriter, r *http.Request) (bson.M, []interface{}) {
	id := UserID(w, r, props)

	var user bson.M
	opts := options.FindOne().SetProjection(bson.M{"_id": 0, "username": 1, "avatar": 1, "trial": 1, "premium": 1})
	props.DB["users"].FindOne(props.Ctx, bson.M{"_id": id}, opts).Decode(&user)

	update.Premium(props, w, r, id, user)

	// Получить список пользователей, которые отправили залогиненому юзеру сообщение хотя бы 1 раз (не прочитанное)
	newMessages, _ := props.DB["messages"].Distinct(props.Ctx, "user", bson.M{"target": id, "viewed": false})
	return user, newMessages
}
