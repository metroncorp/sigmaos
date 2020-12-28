package main

import (
	"log"
	"os"

	"ulambda/fs"
	"ulambda/proc"
)

func main() {
	log.Printf("Running: %v\n", os.Args[0])
	clnt, err := fs.InitFsClient(fs.MakeFsRoot(), os.Args[1:])
	if err != nil {
		log.Fatal("InitFsClient error:", err)
	}

	fd, err := clnt.Open("/")
	if err != nil {
		log.Fatal("Open error:", err)
	}
	defer clnt.Close(fd)
	if buf, err := clnt.Read(fd, 1024); err == nil {
		_, err := clnt.Write(fs.Stdout, buf)
		if err != nil {
			log.Fatal("Write error:", err)
		}
	} else {
		log.Fatal("Read error:", err)
	}
	proc.Exit(clnt)
}
