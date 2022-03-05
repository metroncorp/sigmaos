package snapshot

import (
	"log"

	"ulambda/dir"
	"ulambda/fs"
	"ulambda/inode"
	np "ulambda/ninep"
	"ulambda/protsrv"
)

type Dev struct {
	fs.FsObj
	srv protsrv.FsServer
}

func MakeDev(srv protsrv.FsServer, ctx fs.CtxI, root fs.Dir) *Dev {
	i := inode.MakeInode(ctx, 0, root)
	dev := &Dev{i, srv}
	dir.MkNod(ctx, root, "snapshot", dev)
	return dev
}

func (dev *Dev) Read(ctx fs.CtxI, off np.Toffset, cnt np.Tsize, v np.TQversion) ([]byte, *np.Err) {
	b := dev.srv.Snapshot()
	if len(b) > int(np.MAXGETSET) {
		log.Fatalf("FATAL snapshot too big: %v bytes", len(b))
	}
	return b, nil
}

func (dev *Dev) Write(ctx fs.CtxI, off np.Toffset, b []byte, v np.TQversion) (np.Tsize, *np.Err) {
	log.Printf("Received snapshot of length %v", len(b))
	return np.Tsize(len(b)), nil
}

func (dev *Dev) Snapshot(fn fs.SnapshotF) []byte {
	return dev.FsObj.Snapshot(fn)
}
