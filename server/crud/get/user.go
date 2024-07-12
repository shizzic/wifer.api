package get

import (
	"errors"
	"net/http"
	"strconv"

	unrolled "github.com/unrolled/render"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var render = unrolled.New()

// Получаем id пользователя из куки, для авторизации
func UserID(w http.ResponseWriter, r *http.Request, props Props) (result int) {
	id, err := r.Cookie("id")
	if err != nil {
		render.JSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	result, _ = strconv.Atoi(id.Value)
	return
}

// Получаем весь профиль пользователя
func Profile(id int, props Props) (bson.M, error) {
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
