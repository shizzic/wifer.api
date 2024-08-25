package mongorestore

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
	path := filepath.Join(filename + ".tar.gz")
	file, err := os.Open(path)

	if err != nil {
		log.Fatal("failed to open database archive")
	} else {
		defer os.RemoveAll(props.Conf.PATH + "/" + filename + ".tar.gz")
		defer os.RemoveAll(props.Conf.PATH + "/" + filename)
		defer file.Close()

		if err := extract.Gz(props.Ctx, file, filename, nil); err != nil {
			log.Fatal("failed to extract database archive")
		} else {
			restore(props, "/"+filename+"/"+filename)
		}
	}
}

func restore(props *structs.Props, destination string) {
	cmd := exec.Command("mongorestore", "--uri="+props.Conf.MONGO_CONNECTION_STRING, "-d", "db", props.Conf.PATH+destination)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatal("failed to create an init database\n", err, "\n", out, "\n")
	}
}
