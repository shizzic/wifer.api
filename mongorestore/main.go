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

func extract_archive(props *structs.Props, filename string) {
	archivePath := props.Conf.PATH + "/" + filename + ".tar.gz"
	log.Println("=== EXTRACT DEBUG ===")
	log.Println("Archive path:", archivePath)

	file, err := os.Open(archivePath)
	if err != nil {
		log.Fatal("failed to open database archive: ", err)
	} else {
		info, _ := file.Stat()
		log.Println("Archive size:", info.Size(), "bytes")

		defer os.RemoveAll(archivePath)
		defer os.RemoveAll(props.Conf.PATH + "/" + filename)
		defer file.Close()

		extractPath := props.Conf.PATH + "/" + filename
		log.Println("Extracting to:", extractPath)

		if err := extract.Gz(props.Ctx, file, extractPath, nil); err != nil {
			log.Fatal("failed to extract database archive: ", err)
		} else {
			restore(props, "/"+filename+"/init_dump")
		}
	}
}

func restore(props *structs.Props, destination string) {
	cmd := exec.Command("mongorestore", "--uri="+props.Conf.MONGO_CONNECTION_STRING, "--nsInclude", "db.*", "--drop", fullPath)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatal("failed to create an init database\n", err, "\n", string(out), "\n")
	}
}