package dump

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"wifer/server/structs"

	"github.com/kothar/go-backblaze"
	"github.com/mholt/archiver/v4"
)

func PrepareDB(props *structs.Props) {
	err := exec.Command("mongodump", "--uri="+props.Conf.MONGO_CONNECTION_STRING, "-d", "db", "-o", props.Conf.PATH+"/cron/dump/trash").Run()
	defer os.RemoveAll(props.Conf.PATH + "/cron/dump/trash/db")
	defer os.Remove(props.Conf.PATH + "/cron/dump/trash/db.tar.gz")

	if err == nil {
		Start(props, "/cron/dump/trash/db", "db")
	}
}

// Длеаю архив по выбранному пути
func Start(props *structs.Props, from, name string) {
	// читаю папку фоток для дальнейшей записи
	files, _ := archiver.FilesFromDisk(nil, map[string]string{
		props.Conf.PATH + from: name,
	})
	to, _ := os.Create(props.Conf.PATH + "/cron/dump/trash/" + name + ".tar.gz")
	defer to.Close()
	format := archiver.CompressedArchive{
		Compression: archiver.Gz{},
		Archival:    archiver.Tar{},
	}
	err := format.Archive(context.Background(), to, files)
	if err == nil {
		upload_to_backblaze(props, name)
	}
}

/*
1. Открываю соединение с backblaze
2. Получаю fileID добавленного до этого архива, чтобы удалить его
3. Отправляю новый архив
4. После добавления получаю новый fileID, чтобы воспользоваться им позже
*/
func upload_to_backblaze(props *structs.Props, name string) {
	b2, err := backblaze.NewB2(backblaze.Credentials{
		AccountID:      props.Conf.BACKBLAZE_ID,
		ApplicationKey: props.Conf.BACKBLAZE_KEY,
	})

	if err == nil {
		bucket_name := props.Conf.PRODUCT_NAME + "-bucket"
		b2.CreateBucket(bucket_name, backblaze.AllPrivate)
		bucket, err := b2.Bucket(bucket_name)

		if err == nil {
			file, err := os.Open(props.Conf.PATH + "/cron/dump/trash/" + name + ".txt")
			if err == nil {
				defer file.Close()
				content, _ := os.ReadFile(props.Conf.PATH + "/cron/dump/trash/" + name + ".txt")
				file_id := string(content)
				bucket.DeleteFileVersion(name+".tar.gz", file_id)
			}

			reader, _ := os.Open(props.Conf.PATH + "/cron/dump/trash/" + name + ".tar.gz")
			metadata := make(map[string]string)
			uploaded_file, err := bucket.UploadFile(name+".tar.gz", metadata, reader)
			fmt.Print(err, "\n")
			defer reader.Close()

			// удаляю старый файл и создаю вместо него такой же с актуальным ID
			os.Remove(props.Conf.PATH + "/cron/dump/trash/" + name + ".txt")
			new_file_id, _ := os.Create(props.Conf.PATH + "/cron/dump/trash/" + name + ".txt")
			defer new_file_id.Close()
			new_file_id.WriteString(uploaded_file.ID)
		}
	}
}
