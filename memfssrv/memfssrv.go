package memfssrv

import (
	"sigmaos/dir"
	"sigmaos/fcall"
	"sigmaos/fs"
	"sigmaos/fslib"
	"sigmaos/inode"
	"sigmaos/lockmap"
	"sigmaos/namei"
	"sigmaos/path"
	sp "sigmaos/sigmap"
)

var rootP = path.Path{""}

func (fs *MemFs) Root() fs.Dir {
	return fs.root
}

func (fs *MemFs) FsLib() *fslib.FsLib {
	return fs.fsl
}

// Note: MkDev() sets parent
func (mfs *MemFs) MakeDevInode() *inode.Inode {
	return inode.MakeInode(mfs.ctx, sp.DMDEVICE, nil)
}

func (mfs *MemFs) lookup(path path.Path) (fs.FsObj, *lockmap.PathLock, *fcall.Err) {
	d := mfs.root
	lk := mfs.plt.Acquire(mfs.ctx, rootP)
	if len(path) == 0 {
		return d, lk, nil
	}
	_, lo, lk, _, err := namei.Walk(mfs.plt, mfs.ctx, d, lk, rootP, path, nil)
	if err != nil {
		mfs.plt.Release(mfs.ctx, lk)
		return nil, nil, err
	}
	return lo, lk, nil
}

func (mfs *MemFs) lookupParent(path path.Path) (fs.Dir, *lockmap.PathLock, *fcall.Err) {
	lo, lk, err := mfs.lookup(path)
	if err != nil {
		return nil, nil, err
	}
	d := lo.(fs.Dir)
	return d, lk, nil
}

func (mfs *MemFs) MkDev(pn string, dev fs.Inode) *fcall.Err {
	path := path.Split(pn)
	d, lk, err := mfs.lookupParent(path.Dir())
	if err != nil {
		return err
	}
	defer mfs.plt.Release(mfs.ctx, lk)
	// i := inode.MakeInode(mfs.ctx, sp.DMDEVICE, d)
	dev.SetParent(d)
	return dir.MkNod(mfs.ctx, d, path.Base(), dev)
}

func (mfs *MemFs) MkNod(pn string, i fs.Inode) *fcall.Err {
	path := path.Split(pn)
	d, lk, err := mfs.lookupParent(path.Dir())
	if err != nil {
		return err
	}
	defer mfs.plt.Release(mfs.ctx, lk)
	return dir.MkNod(mfs.ctx, d, path.Base(), i)
}

func (mfs *MemFs) Create(pn string, p sp.Tperm, m sp.Tmode) (fs.FsObj, *fcall.Err) {
	path := path.Split(pn)
	d, lk, err := mfs.lookupParent(path.Dir())
	if err != nil {
		return nil, err
	}
	defer mfs.plt.Release(mfs.ctx, lk)
	return d.Create(mfs.ctx, path.Base(), p, m)
}

func (mfs *MemFs) Remove(pn string) *fcall.Err {
	path := path.Split(pn)
	d, lk, err := mfs.lookupParent(path.Dir())
	if err != nil {
		return err
	}
	defer mfs.plt.Release(mfs.ctx, lk)
	return d.Remove(mfs.ctx, path.Base())
}

func (mfs *MemFs) Open(pn string, m sp.Tmode) (fs.FsObj, *fcall.Err) {
	path := path.Split(pn)
	lo, lk, err := mfs.lookup(path)
	if err != nil {
		return nil, err
	}
	mfs.plt.Release(mfs.ctx, lk)
	return lo, nil
}
