package mongorestore

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"wifer/server/structs"

	"github.com/codeclysm/extract/v4"
)

// Начать процесс создания базы данных, если ее нету по какой то причине
func Start(props *structs.Props) error {
	_, err := os.Stat(props.Conf.PATH + "/init_dump")

	if osname := runtime.GOOS; osname == "windows" && err == nil {
		restore(props, "/init_dump")
	} else {
		download_backblaze_dump(props, osname)
	}

	return nil
}

// реархивировать скаченный dump
func extract_archive(props *structs.Props, filename string) {
	file, err := os.Open(filename + ".tar.gz")

	if err != nil {
		log.Fatal("failed to open database archive")
	} else {
		defer os.RemoveAll(filename + ".tar.gz")
		defer os.RemoveAll(filename)
		defer file.Close()

		if err := extract.Gz(props.Ctx, file, filename, nil); err != nil {
			log.Fatal("failed to extract database archive")
		} else {
			restore(props, "/"+filename+"/"+filename)
		}
	}
}

func restore(props *structs.Props, destination string) {
	exec.Command("ls").Run()
	err := exec.Command("mongorestore", "--uri="+props.Conf.MONGO_CONNECTION_STRING, "-d", "db", props.Conf.PATH+destination).Run()

	if err != nil {
		log.Fatal(err.Error(), "\n", "failed to create an init database")
	}
}
