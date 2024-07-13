package routes

import (
	"io"
	"net/http"
	"os"
)

func image(props *Props) {
	props.R.Post("/upload-image", func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(20 << 20)                                  // указываем максимальный размер файла, 20мб
		file, _, _ := r.FormFile("file")                                // извлекаю файл с именем "file"
		F, _ := os.OpenFile("image.jpg", os.O_WRONLY|os.O_CREATE, 0666) // создаю пустышку для будущего файла в текущем каталоге
		io.Copy(F, file)                                                // Копирую файл в переменную F, то есть в каталог скрипта

		// files := r.MultipartForm.File["file[]"]
		// for i, header := range files {
		// 	file, _ := header.Open()
		// 	fmt.Print(i)
		// 	fmt.Print("\n")
		// }
	})
}
