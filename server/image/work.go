package image

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"wifer/server/structs"
)

// Изменить расположение фотки
func ChangeDir(props *structs.Props, data *Images) (err error) {
	var from string
	var to string

	if data.Into == "" {
		// Вернуть ошибку, если при переносе аватарки, ничем ее заменить
		if data.Count == 1 {
			err = errors.New("avatar_must_be")
			return
		}

		// Указать, что аватара больше нет
		data.IsAvatar = false
		from = data.Avatar
	} else {
		from = data.Path + "/" + data.Into + "/" + data.Filename + ".webp"
	}
	to = data.Path + "/" + data.NewDir + "/new.webp"
	os.Rename(from, to)
	restoreAvatar(data)
	rename(props, data)
	return
}

// Поставить аву из имеющихся фоток, если ее по какой то причине больше нету
// В приоритете публичные папка
func restoreAvatar(data *Images) {
	if !data.IsAvatar {
		dirs := [2]string{"public", "private"}

	out:
		for _, dir := range dirs {
			path := data.Path + "/" + dir
			files, _ := os.ReadDir(path)

			for _, image := range files {
				if image.Name() != "new.webp" {
					from := path + "/" + image.Name()
					to := data.Avatar
					fmt.Print("\n", from, "\n", to, "\n")
					os.Rename(from, to)
					break out
				}
			}
		}
	}
}

// Переименную фотки в каждой папке с 1 и до конца
func rename(props *structs.Props, data *Images) {
	dirs := [2]string{"public", "private"}

	for _, dir := range dirs {
		path := data.Path + "/" + dir
		files, _ := os.ReadDir(path)

		for index, image := range files {
			new_name := strconv.Itoa(index+1) + ".webp"

			if image.Name() != new_name {
				from := path + "/" + image.Name()
				to := path + "/" + new_name
				os.Rename(from, to)
			}
		}
	}

	count(data)
	updateCounts(props, data)
}
