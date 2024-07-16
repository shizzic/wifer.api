package image

import (
	"errors"
	"os"
	"wifer/server/structs"
)

// Изменить расположение фотки
func ChangeDir(props *structs.Props, data *Images) (err error) {
	// Вернуть ошибку, если при переносе аватарки, ничем ее заменить
	if data.Into == "" && data.Count == 1 {
		err = errors.New("avatar_must_be")
		return
	}

	from := data.Path + "/" + data.Into + data.Filename // filename -> /{name}.webp
	to := data.Path + "/" + data.NewDir + "/new.webp"
	move(from, to)

	count(data)
	restoreAvatar(data)
	rename(props, data)
	return
}

// Свапнуть аватарку и выбранную фотку
func ReplaceAvatar(props *structs.Props, data *Images) {
	move(data.Path+"/"+data.Into+"/"+data.Filename+".webp", data.Path+"/new_avatar.webp")
	move(data.Avatar, data.Path+"/"+data.Into+"/"+data.Filename+".webp")
	move(data.Path+"/new_avatar.webp", data.Avatar)

	count(data)
	rename(props, data)
}

// Удаляю фото
// Если это была аватарка, заменяю ее на первую попавшуюся фотку
func DeleteImage(props *structs.Props, data *Images) {
	if err := os.Remove(data.Path + data.Into + "/" + data.Filename + ".webp"); err != nil {
		DeleteImage(props, data)
	}

	count(data)
	restoreAvatar(data)
	rename(props, data)
}
