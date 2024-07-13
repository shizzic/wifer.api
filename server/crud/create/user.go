package create

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"wifer/server/auth"
	"wifer/server/mail"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Signin(props *auth.Props, w http.ResponseWriter, email string, isApi bool) (int, error) {
	if !auth.IsEmailValid(email) {
		return 0, errors.New("1")
	}

	code := auth.MakeCode()

	var user bson.M
	opts := options.FindOne().SetProjection(bson.M{"_id": 1, "username": 1, "email": 1, "status": 1, "active": 1})
	there_is_no_such_user := props.DB["users"].FindOne(props.Ctx, bson.M{"email": email}, opts).Decode(&user)

	// Если такого пользователя нет, регистрирую
	if there_is_no_such_user != nil {
		// Получаем id последнего юзера
		var last bson.M
		opts = options.FindOne().SetProjection(bson.M{"_id": 1}).SetSort(bson.D{{Key: "_id", Value: -1}})
		props.DB["users"].FindOne(props.Ctx, bson.M{}, opts).Decode(&last)

		id := 1
		if last["_id"] != nil {
			id = int(last["_id"].(int32)) + 1
		}
		date := time.Now().Unix()

		ObjectId, err := props.DB["users"].InsertOne(props.Ctx, bson.D{
			{Key: "_id", Value: id},
			{Key: "username", Value: strconv.Itoa(id)},
			{Key: "email", Value: email},
			{Key: "title", Value: ""},
			{Key: "about", Value: ""},
			{Key: "is_about", Value: false},
			{Key: "sex", Value: 0},
			{Key: "age", Value: 0},
			{Key: "body", Value: 0},
			{Key: "height", Value: 0},
			{Key: "weight", Value: 0},
			{Key: "smokes", Value: 0},
			{Key: "drinks", Value: 0},
			{Key: "ethnicity", Value: 0},
			{Key: "search", Value: []int{}},
			{Key: "prefer", Value: 0},
			{Key: "income", Value: 0},
			{Key: "children", Value: 0},
			{Key: "industry", Value: 0},
			{Key: "country_id", Value: 0},
			{Key: "city_id", Value: 0},
			{Key: "premium", Value: int64(0)},
			{Key: "trial", Value: false},
			{Key: "status", Value: isApi},
			{Key: "active", Value: isApi},
			{Key: "created_at", Value: date},
			{Key: "last_time", Value: date},
			{Key: "online", Value: false},
			{Key: "avatar", Value: false},
			{Key: "public", Value: 0},
			{Key: "private", Value: 0},
			{Key: "images", Value: 0},
		})

		if err != nil {
			return 0, errors.New("3")
		}

		if isApi {
			id := strconv.Itoa(int(ObjectId.InsertedID.(int32)))
			auth.MakeCookies(props, id, id, 86400*120, w)
			return int(ObjectId.InsertedID.(int32)), nil
		} else {
			if _, err := props.DB["ensure"].InsertOne(props.Ctx, bson.D{
				{Key: "_id", Value: ObjectId.InsertedID},
				{Key: "code", Value: code},
			}); err != nil {
				// Delete new user, because code wasn't added
				props.DB["users"].DeleteOne(props.Ctx, bson.M{"_id": int(ObjectId.InsertedID.(int32))})
				return 0, errors.New("3")
			}

			if err := mail.SendCode(props, email, code, strconv.Itoa(int(ObjectId.InsertedID.(int32)))); err != nil {
				return 0, errors.New("2")
			}
		}
	} else {
		if !user["status"].(bool) {
			return 0, errors.New("4")
		}

		if !user["active"].(bool) {
			props.DB["users"].UpdateOne(props.Ctx, bson.M{"_id": user["_id"].(int32)}, bson.D{{Key: "$set", Value: bson.D{{Key: "active", Value: true}}}})
		}

		if isApi {
			auth.MakeCookies(props, strconv.Itoa(int(user["_id"].(int32))), user["username"].(string), 86400*120, w)
			return int(user["_id"].(int32)), nil
		} else {
			props.DB["ensure"].DeleteOne(props.Ctx, bson.M{"_id": user["_id"]})
			props.DB["ensure"].InsertOne(props.Ctx, bson.D{
				{Key: "_id", Value: user["_id"]},
				{Key: "code", Value: code},
			})

			if err := mail.SendCode(props, email, code, strconv.Itoa(int(user["_id"].(int32)))); err != nil {
				return 0, errors.New("2")
			}
		}
	}

	return 0, nil
}
