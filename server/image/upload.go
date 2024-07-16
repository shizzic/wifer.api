package image

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"wifer/server/structs"

	"github.com/nickalie/go-webpbin"
)

type Images = structs.Images

func Upload(props *structs.Props, r *http.Request, data *Images) (overcount int8, err error) {
	// Проверяю что кол-во фоток не превысит 20 после добавления
	files := r.MultipartForm.File["files[]"]
	new_files_count := int8(len(files))
	if (data.Count + new_files_count) > 20 {
		overcount = data.Count + new_files_count - 20
		err = errors.New("max_image")
		return
	}

	for _, header := range files {
		// 20 мб
		if header.Size <= (20 << 20) {
			data.Count++
			if data.IsAvatar {
				if data.Into == "public" {
					data.CountPublic++
				} else {
					data.CountPrivate++
				}

				data.Output = data.FullPath + strconv.FormatInt(int64(data.Count-1), 10) + ".webp"
			} else {
				data.Output = data.Avatar
			}
			file, _ := header.Open()
			defer file.Close() // закрытие
			save(props, data, file)
		}
	}

	rename(props, data)
	return
}

// Сохранить пришедшее фото
func save(props *structs.Props, data *Images, file multipart.File) {
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, file)
	output := data.Output

	// доавить аву, если ее нет
	if !data.IsAvatar {
		output = data.Avatar
		data.IsAvatar = true
	}

	os.WriteFile(output, buf.Bytes(), 0666)
	convertToWEBP(output)
	updateCounts(props, data)
}

// Конвертировать сохраненную фотку в webp
func convertToWEBP(output string) {
	webpbin.NewCWebP().
		Quality(100).
		InputFile(output).
		OutputFile(output).
		Run()
}
