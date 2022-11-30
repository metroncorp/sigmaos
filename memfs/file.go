package memfs

import (
	//"time"

	"sigmaos/file"
	"sigmaos/fs"
	np "sigmaos/sigmap"
    "sigmaos/fcall"
)

type File struct {
	fs.Inode
	*file.File
}

func MakeFile(i fs.Inode) *File {
	f := &File{}
	f.Inode = i
	f.File = file.MakeFile()
	return f
}

func (f *File) Size() (np.Tlength, *fcall.Err) {
	return f.File.Size()
}

func (f *File) Stat(ctx fs.CtxI) (*np.Stat, *fcall.Err) {
	st, err := f.Inode.Stat(ctx)
	if err != nil {
		return nil, err
	}
	st.Length, _ = f.Size()
	return st, nil
}

func (f *File) Snapshot(fn fs.SnapshotF) []byte {
	return makeFileSnapshot(f)
}

func RestoreFile(fn fs.RestoreF, b []byte) fs.Inode {
	return restoreFile(fn, b)
}
