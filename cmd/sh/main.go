package main

import (
	"log"
	"os"
	"strings"

	"ulambda/fs"
	"ulambda/proc"
)

func main() {
	log.Printf("Running: %v\n", os.Args[0])
	clnt, err := fs.InitFsClient(fs.MakeFsRoot(), os.Args[1:])
	if err != nil {
		log.Fatal("InitFsClient error:", err)
	}
	for {
		b := []byte("λ ")
		_, err := clnt.Write(fs.Stdout, b)
		if err != nil {
			log.Fatal("Write error:", err)
		}
		b, err = clnt.Read(fs.Stdin, 1024)
		if err != nil {
			log.Fatal("Read error:", err)
		}
		cmd := strings.TrimSuffix(string(b), "\n")
		child, err := proc.Spawn(clnt, cmd, clnt.Lsof())
		if err != nil {
			log.Fatal("Spawn error:", err)
		}
		status, err := proc.Wait(clnt, child)
		if err != nil {
			log.Fatal("Wait error:", err)
		}
		log.Printf("Wait: %v\n", string(status))
	}
	proc.Exit(clnt)
}
