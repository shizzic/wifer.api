package mongorestore

import (
	"bytes"
	"log"
	"os"
	"wifer/server/structs"

	"github.com/kothar/go-backblaze"
)

func download_backblaze_dump(props *structs.Props, osname string) {
	log.Println("=== BACKBLAZE DEBUG ===")
	log.Println("OS:", osname)
	log.Println("PATH:", props.Conf.PATH)

	if b2, err := backblaze.NewB2(backblaze.Credentials{
		AccountID:      props.Conf.BACKBLAZE_ID,
		ApplicationKey: props.Conf.BACKBLAZE_KEY,
	}); err != nil {
		log.Fatal("could not connect to backblaze: ", err)
	} else {
		bucket_name := props.Conf.PRODUCT_NAME + "-bucket"
		log.Println("Bucket name:", bucket_name)
		b2.CreateBucket(bucket_name, backblaze.AllPrivate)

		if bucket, err := b2.Bucket(bucket_name); err != nil {
			log.Fatal("backblaze's bucket was not created for some reason: ", err)
		} else {
			filename := "init_dump"

			if list, err := bucket.ListFileNamesWithPrefix("db.tar.gz", 1, "db.tar.gz", ""); err != nil {
				log.Fatal("failed to list filenames in backblaze storage: ", err)
			} else {
				log.Println("Files found in backblaze:", len(list.Files))
				for _, f := range list.Files {
					log.Println("File:", f.Name, f.Size)
				}

				if osname != "windows" && len(list.Files) > 0 {
					filename = "db"
				}
				log.Println("Downloading filename:", filename)

				if file, reader, err := bucket.DownloadFileByName(filename + ".tar.gz"); err != nil {
					log.Fatal("failed to download database: ", err)
				} else {
					log.Println("Downloaded file ID:", file.ID, "Size:", file.Size)

					os.Remove(props.Conf.PATH + "/cron/dump/trash/db.txt")
					new_file_id, _ := os.Create(props.Conf.PATH + "/cron/dump/trash/db.txt")
					defer new_file_id.Close()
					new_file_id.WriteString(file.ID)

					defer reader.Close()
					buf := new(bytes.Buffer)

					if _, err := buf.ReadFrom(reader); err != nil {
						log.Fatal("could not read file into bytes from backblaze: ", err)
					} else {
						log.Println("Buffer size:", buf.Len(), "bytes")

						if err := os.WriteFile(props.Conf.PATH+"/"+filename+".tar.gz", buf.Bytes(), 0755); err != nil {
							log.Fatal("failed to save downloaded dump: ", err)
						} else {
							log.Println("Saved archive, starting extraction")
							extract_archive(props, filename)
						}
					}
				}
			}
		}
	}
}