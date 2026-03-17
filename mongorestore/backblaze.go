package mongorestore

import (
	"bytes"
	"log"
	"os"
	"wifer/server/structs"

	"github.com/kothar/go-backblaze"
)

func download_backblaze_dump(props *structs.Props, osname string) {
	if b2, err := backblaze.NewB2(backblaze.Credentials{
		AccountID:      props.Conf.BACKBLAZE_ID,
		ApplicationKey: props.Conf.BACKBLAZE_KEY,
	}); err != nil {
		log.Fatal("could not connect to backblaze: ", err)
	} else {
		bucket_name := props.Conf.PRODUCT_NAME + "-bucket"
		b2.CreateBucket(bucket_name, backblaze.AllPrivate)

		if bucket, err := b2.Bucket(bucket_name); err != nil {
			log.Fatal("backblaze's bucket was not created for some reason: ", err)
		} else {
			filename := "init_dump"

			if list, err := bucket.ListFileNamesWithPrefix("db.tar.gz", 1, "db.tar.gz", ""); err != nil {
				log.Fatal("failed to list filenames in backblaze storage: ", err)
			} else {
				if osname != "windows" && len(list.Files) > 0 {
					filename = "db"
				}

				if file, reader, err := bucket.DownloadFileByName(filename + ".tar.gz"); err != nil {
					log.Fatal("failed to download database: ", err)
				} else {
					os.Remove(props.Conf.PATH + "/cron/dump/trash/db.txt")
					new_file_id, _ := os.Create(props.Conf.PATH + "/cron/dump/trash/db.txt")
					defer new_file_id.Close()
					new_file_id.WriteString(file.ID)

					defer reader.Close()
					buf := new(bytes.Buffer)

					if _, err := buf.ReadFrom(reader); err != nil {
						log.Fatal("could not read file into bytes from backblaze: ", err)
					} else {

						if err := os.WriteFile(props.Conf.PATH+"/"+filename+".tar.gz", buf.Bytes(), 0755); err != nil {
							log.Fatal("failed to save downloaded dump: ", err)
						} else {
							extract_archive(props, filename)
						}
					}
				}
			}
		}
	}
}