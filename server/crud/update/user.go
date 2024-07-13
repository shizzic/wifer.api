package update

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	"wifer/server/auth"
	"wifer/server/structs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Props = structs.Props
type User = structs.User

func Change(props *Props, r *http.Request, w http.ResponseWriter, data *User, id int) error {
	if !auth.IsUsernameValid(data.Username) {
		return errors.New("0")
	}
	if !auth.IsAboutValid(data.About) {
		return errors.New("0")
	}
	if !auth.IsTitleValid(data.Title) {
		return errors.New("0")
	}
	if !auth.IsSexValid(data.Sex) {
		return errors.New("0")
	}
	if !auth.IsAgeValid(data.Age) {
		return errors.New("0")
	}
	if !auth.IsHeightValid(data.Height) {
		return errors.New("0")
	}
	if !auth.IsWeightValid(data.Weight) {
		return errors.New("0")
	}
	if !auth.IsSmokeValid(data.Smokes) {
		return errors.New("0")
	}
	if !auth.IsDrinkValid(data.Drinks) {
		return errors.New("0")
	}
	if !auth.IsEthnicityValid(data.Ethnicity) {
		return errors.New("0")
	}
	if !auth.IsBodyValid(data.Body) {
		return errors.New("0")
	}
	if !auth.IsIncomeValid(data.Income) {
		return errors.New("0")
	}
	if !auth.IsIndustryValid(data.Industry) {
		return errors.New("0")
	}
	if !auth.IsPreferValid(data.Prefer) {
		return errors.New("0")
	}
	if !auth.IsChildrenValid(data.Children) {
		return errors.New("0")
	}

	data.Username = strings.TrimSpace(data.Username)
	data.Title = strings.TrimSpace(data.Title)
	data.About = strings.TrimSpace(data.About)

	oldUsername, _ := r.Cookie("username")
	if oldUsername.Value != data.Username {
		if available := isUsernameAvailable(props, data.Username); !available {
			return errors.New("1")
		}
	}

	isAbout := true
	if len(data.About) == 0 {
		isAbout = false
	}

	if _, err := props.DB["users"].UpdateOne(props.Ctx, bson.M{"_id": id}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "username", Value: data.Username}}},
		{Key: "$set", Value: bson.D{{Key: "title", Value: data.Title}}},
		{Key: "$set", Value: bson.D{{Key: "about", Value: data.About}}},
		{Key: "$set", Value: bson.D{{Key: "is_about", Value: isAbout}}},
		{Key: "$set", Value: bson.D{{Key: "sex", Value: data.Sex}}},
		{Key: "$set", Value: bson.D{{Key: "age", Value: data.Age}}},
		{Key: "$set", Value: bson.D{{Key: "body", Value: data.Body}}},
		{Key: "$set", Value: bson.D{{Key: "weight", Value: data.Weight}}},
		{Key: "$set", Value: bson.D{{Key: "height", Value: data.Height}}},
		{Key: "$set", Value: bson.D{{Key: "smokes", Value: data.Smokes}}},
		{Key: "$set", Value: bson.D{{Key: "drinks", Value: data.Drinks}}},
		{Key: "$set", Value: bson.D{{Key: "ethnicity", Value: data.Ethnicity}}},
		{Key: "$set", Value: bson.D{{Key: "search", Value: data.Search}}},
		{Key: "$set", Value: bson.D{{Key: "prefer", Value: data.Prefer}}},
		{Key: "$set", Value: bson.D{{Key: "income", Value: data.Income}}},
		{Key: "$set", Value: bson.D{{Key: "children", Value: data.Children}}},
		{Key: "$set", Value: bson.D{{Key: "industry", Value: data.Industry}}},
		{Key: "$set", Value: bson.D{{Key: "country_id", Value: data.Country}}},
		{Key: "$set", Value: bson.D{{Key: "city_id", Value: data.City}}},
	}); err != nil {
		return errors.New("2")
	}

	strID := strconv.Itoa(id)
	auth.MakeCookies(props, strID, data.Username, 86400*120, w)
	return nil
}

func isUsernameAvailable(props *Props, username string) bool {
	var data bson.M
	opts := options.FindOne().SetProjection(bson.M{"username": 1})
	if err := props.DB["users"].FindOne(props.Ctx, bson.M{"username": username}, opts).Decode(&data); err == nil {
		return false
	}

	return true
}

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
