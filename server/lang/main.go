package lang

import gt "github.com/bas24/googletranslatefree"

// Перевод строки
func Translate(text, lang string) (string, error) {
	result, err := gt.Translate(text, "auto", lang)

	if err != nil {
		return "", err
	} else {
		return result, nil
	}
}
