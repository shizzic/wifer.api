package get

import (
	"regexp"
	"strings"
	"wifer/server/structs"

	gt "github.com/bas24/googletranslatefree"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Countries(props *structs.Props, locale string) []bson.M {
	var data []bson.M
	cursor, _ := props.DB["countries"].Find(props.Ctx, bson.M{})
	cursor.All(props.Ctx, &data)
	data = translate(data, locale)

	return data
}

func Cities(props *structs.Props, country_id int, locale string) []bson.M {
	var data []bson.M
	opts := options.Find().SetProjection(bson.M{"_id": 1, "title": 1})
	cursor, _ := props.DB["cities"].Find(props.Ctx, bson.M{"country_id": country_id}, opts)
	cursor.All(props.Ctx, &data)
	data = translate(data, locale)

	return data
}

// Перевод локаций
func translate(data []bson.M, locale string) []bson.M {
	from := 0
	divider := ";"
	slices := split(data)

	for _, slice := range slices {
		str := ""

		for _, value := range slice {
			str += divider + value["title"].(string) + divider
		}

		str = strings.Trim(str, divider)
		translated, _ := gt.Translate(str, "en", locale)
		// translated = strings.Replace(translated, " ", "", -1) 	// удаление вообще всех пробелов в строке
		trimmed := removeWhiteSpaces(translated)
		splitted := strings.Split(trimmed, divider+divider)

		for _, value := range splitted {
			data[from]["title"] = value
			from++
		}
	}

	return data
}

// разделить локации на массив массивов локаций для перевода
func split(data []bson.M) (result [][]bson.M) {
	var numberOfChunks int
	if len(data) <= 150 {
		numberOfChunks = 1
	} else {
		numberOfChunks = len(data) / 150 // в каждом подмассиве 150 элементов
	}

	for i := 0; i < numberOfChunks; i++ {
		min := (i * len(data) / numberOfChunks)
		max := ((i + 1) * len(data)) / numberOfChunks

		result = append(result, data[min:max])
	}

	return
}

// Удаляю лишнии пробелы в переведенной строке
func removeWhiteSpaces(str string) (result string) {
	re := regexp.MustCompile(`;[\s,]*;`)
	result = re.ReplaceAllString(str, `;;`)
	re = regexp.MustCompile(`[\s,]*;;[\s,]*`)
	result = re.ReplaceAllString(result, `;;`)
	return
}
