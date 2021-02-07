package fslambda

import (
  "io/ioutil"

  "ulambda/fslib"
  db "ulambda/debug"
)

type Downloader struct {
  pid  string
  src  string
  dest string
  *fslib.FsLib
}

func MakeDownloader(args []string, debug bool) (*Downloader, error) {
  db.DPrintf("Downloader: %v\n", args)
  down := &Downloader{}
  down.pid = args[0]
  down.src = args[1]
  down.dest = args[2]
  fls := fslib.MakeFsLib("downloader")
  down.FsLib = fls
  db.SetDebug(debug)
  down.Started(down.pid)
  return down, nil
}

func (down *Downloader) Work() {
  db.DPrintf("Downloading [%v] to [%v]\n", down.src, down.dest);
  contents, err := down.ReadFile(down.src)
  if err != nil {
    db.DPrintf("Read download file error [%v]: %v\n", down.src, err)
  }
  err = ioutil.WriteFile(down.dest, contents, 0666)
  if err != nil {
    db.DPrintf("Couldn't download file [%v]: %v\n", down.dest, err)
  }
}

func (down *Downloader) Exit() {
  down.Exiting(down.pid)
}
