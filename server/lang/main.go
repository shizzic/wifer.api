package lang

import (
	"wifer/server/structs"

	gt "github.com/bas24/googletranslatefree"
)

type Translate = structs.Translate

// Перевод строки
func TranslateText(props *Translate) (string, error) {
	result, err := gt.Translate(props.Text, "auto", props.Lang)

	if err != nil {
		return "", err
	} else {
		return result, nil
	}
}
