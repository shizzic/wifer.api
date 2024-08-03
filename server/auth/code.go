package auth

import (
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Создать код для валидации любого действия, требующего подтверждения
func MakeCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 6)
	for i := range b {
		b[i] = nums[r.Int63()%int64(len(nums))]
	}

	return string(b)
}

// Проверить код для валидации
func CheckCode(props *Props, id int, code string, w http.ResponseWriter) error {
	if !isCode(code) {
		return errors.New("invalid_code")
	}

	var data bson.M
	opts := options.FindOne().SetProjection(bson.M{"_id": 1})

	if err := props.DB["ensure"].FindOne(props.Ctx, bson.M{"_id": id, "code": code}, opts).Decode(&data); err != nil {
		return errors.New("code_not_found")
	}

	// Delete document in ensure collection, if given code was valid
	props.DB["ensure"].DeleteOne(props.Ctx, bson.M{"_id": id, "code": code})

	var user bson.M
	opt := options.FindOne().SetProjection(bson.M{"username": 1})

	if err := props.DB["users"].FindOne(props.Ctx, bson.M{"_id": id}, opt).Decode(&user); err != nil {
		return errors.New("user_not_found")
	}

	props.DB["users"].UpdateOne(props.Ctx, bson.M{"_id": id}, bson.D{{Key: "$set", Value: bson.D{
		{Key: "status", Value: true},
		{Key: "active", Value: true},
	}}})

	MakeCookies(props, strconv.Itoa(id), user["username"].(string), 86400*120, w)

	return nil
}
