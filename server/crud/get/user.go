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
func UserID(w http.ResponseWriter, r *http.Request, props *structs.Props) (result int) {
	id, err := r.Cookie("id")
	if err != nil {
		render.JSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	result, _ = strconv.Atoi(id.Value)
	return
}

// Получаем весь профиль пользователя
func Profile(id int, props *structs.Props) (bson.M, error) {
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
		"images":     1,
	})

	if err := props.DB["users"].FindOne(props.Ctx, bson.M{"_id": id, "status": true}, opts).Decode(&user); err != nil {
		return user, errors.New("no_such_user_or_banned")
	}

	return user, nil
}

func UserEmailByApi(props *structs.Props, data *Signin) (email string, err error) {
	switch data.Method {
	case "Google":
		email, err = auth.IsGoogle(data.ID, data.Token)
	case "Yandex":
		email, err = auth.IsYandex(props, data.Token)
	case "Mail":
		email, err = auth.IsMail(props, data)
	}
	return
}

// Получаю пользователя и список айдишников юзеров, которые ему отписали (не прочитанные сообщения)
func UserMainInfo(props *structs.Props, w http.ResponseWriter, r *http.Request) (bson.M, bool) {
	id := UserID(w, r, props)

	var user bson.M
	opts := options.FindOne().SetProjection(bson.M{"_id": 0, "username": 1, "avatar": 1, "trial": 1, "premium": 1})
	props.DB["users"].FindOne(props.Ctx, bson.M{"_id": id}, opts).Decode(&user)

	premium := update.Premium(props, w, r, id, user)
	return user, premium
}

// Фильтр полей для поиска пользователей
func PrepareFilter(data *structs.Template) bson.M {
	filter := bson.M{
		"age":      bson.M{"$gte": data.AgeMin, "$lte": data.AgeMax},
		"height":   bson.M{"$gte": data.HeightMin, "$lte": data.HeightMax},
		"weight":   bson.M{"$gte": data.WeightMin, "$lte": data.WeightMax},
		"children": bson.M{"$gte": data.ChildrenMin, "$lte": data.ChildrenMax},
		"images":   bson.M{"$gte": data.ImagesMin, "$lte": data.ImagesMax},
	}

	if len(data.Body) > 0 {
		filter["body"] = bson.M{"$in": data.Body}
	}

	if len(data.Sex) > 0 {
		filter["sex"] = bson.M{"$in": data.Sex}
	}

	if len(data.Smokes) > 0 {
		filter["smokes"] = bson.M{"$in": data.Smokes}
	}

	if len(data.Drinks) > 0 {
		filter["drinks"] = bson.M{"$in": data.Drinks}
	}

	if len(data.Ethnicity) > 0 {
		filter["ethnicity"] = bson.M{"$in": data.Ethnicity}
	}

	if len(data.Search) > 0 {
		filter["search"] = bson.M{"$in": data.Search}
	}

	if len(data.Income) > 0 {
		filter["income"] = bson.M{"$in": data.Income}
	}

	if len(data.Industry) > 0 {
		filter["industry"] = bson.M{"$in": data.Industry}
	}

	if len(data.Premium) > 0 {
		filter["premium"] = bson.M{"$in": data.Premium}
	}

	if len(data.Country) > 0 {
		filter["country_id"] = bson.M{"$in": data.Country}
	}

	if len(data.City) > 0 {
		filter["city_id"] = bson.M{"$in": data.City}
	}

	if data.IsAbout {
		filter["is_about"] = true
	}

	if data.Avatar {
		filter["avatar"] = true
	}

	if data.Text != "" {
		filter["$text"] = bson.M{"$search": data.Text}
	}

	filter["status"] = true
	filter["active"] = true

	return filter
}

// Посчитать кол-во пользователей по заданному фильтру
func CountUsersByFilter(props *structs.Props, filter *bson.M) int64 {
	count, err := props.DB["users"].CountDocuments(props.Ctx, filter)

	if err != nil {
		return 0
	} else {
		return count
	}
}

// Получаю пользователей по выбранному (одному) фильтру
func UsersByFilter(props *structs.Props, data *structs.Template, filter *bson.M) (list []bson.M) {
	opts := options.Find().SetProjection(bson.M{
		"username":   1,
		"title":      1,
		"age":        1,
		"weight":     1,
		"height":     1,
		"body":       1,
		"ethnicity":  1,
		"public":     1,
		"private":    1,
		"avatar":     1,
		"premium":    1,
		"country_id": 1,
		"city_id":    1,
		"online":     1,
		"is_about":   1,
	}).
		SetSort(bson.D{
			{Key: "premium", Value: -1},
			{Key: data.Sort, Value: -1},
			{Key: "_id", Value: 1},
		}).
		SetLimit(data.Limit).
		SetSkip(data.Skip)

	cursor, _ := props.DB["users"].Find(props.Ctx, filter, opts)
	cursor.All(props.Ctx, &list)
	return
}

func CheckUsernameAvailability(props *structs.Props, username string) bool {
	if auth.IsUsernameValid(username) {
		var data bson.M
		opts := options.FindOne().SetProjection(bson.M{"username": 1})
		if err := props.DB["users"].FindOne(props.Ctx, bson.M{"username": username}, opts).Decode(&data); err == nil {
			return false
		}

		return true
	}

	return false
}
