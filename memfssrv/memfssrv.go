// Servers use package memfsssrv to create an in-memory file server.
// memfsssrv uses protsrv to handle client sigmaP requests. The server
// itself can also create, remove, etc.  files in its file system,
// which directly call into protsrv.
package memfssrv

import (
	"path/filepath"

	"sigmaos/auth"
	"sigmaos/ctx"
	db "sigmaos/debug"
	"sigmaos/dir"
	"sigmaos/fs"
	"sigmaos/inode"
	"sigmaos/lockmap"
	"sigmaos/namei"
	"sigmaos/path"
	"sigmaos/portclnt"
	"sigmaos/protsrv"
	"sigmaos/serr"
	"sigmaos/sigmaclnt"
	sp "sigmaos/sigmap"
	"sigmaos/sigmapsrv"
	"sigmaos/syncmap"
)

var rootP = path.Path{""}

type MemFs struct {
	*sigmapsrv.SigmaPSrv
	ctx   fs.CtxI // server context
	plt   *lockmap.PathLockTable
	ps    *protsrv.ProtSrv
	roots *roots
	sc    *sigmaclnt.SigmaClnt
	amgr  auth.AuthMgr
	pc    *portclnt.PortClnt
	pi    portclnt.PortInfo
	pn    string
}

type roots struct {
	roots *syncmap.SyncMap[string, sp.Tfid]
	next  sp.Tfid
}

func newRoots() *roots {
	r := &roots{roots: syncmap.NewSyncMap[string, sp.Tfid]()}
	return r
}

func (rts *roots) lookupAlloc(pn string) (sp.Tfid, bool) {
	fid, ok := rts.roots.AllocNew(pn, func(k string) sp.Tfid {
		db.DPrintf(db.MEMFSSRV, "lookupAlloc: new root %q fid %d\n", pn, rts.next+1)
		rts.next += 1
		return rts.next
	})
	db.DPrintf(db.MEMFSSRV, "lookupAlloc: root %q fid %d\n", pn, fid)
	return fid, ok
}

func NewMemFsSrv(pn string, srv *sigmapsrv.SigmaPSrv, sc *sigmaclnt.SigmaClnt, amgr auth.AuthMgr, fencefs fs.Dir) *MemFs {
	mfs := &MemFs{
		SigmaPSrv: srv,
		ctx:       ctx.NewCtx(sc.ProcEnv().GetPrincipal(), nil, 0, sp.NoClntId, nil, fencefs),
		plt:       srv.PathLockTable(),
		sc:        sc,
		amgr:      amgr,
		pn:        pn,
		ps:        protsrv.NewProtSrv(srv.ProtSrvState, 0, srv.GetRootCtx),
		roots:     newRoots(),
	}
	return mfs
}

func (mfs *MemFs) SigmaClnt() *sigmaclnt.SigmaClnt {
	return mfs.sc
}

// Note: NewDev() sets parent
func (mfs *MemFs) NewDevInode() *inode.Inode {
	return inode.NewInode(mfs.ctx, sp.DMDEVICE, nil)
}

// Returns fid for root and remaining path from root
func (mfs *MemFs) rootFid(pn string) (sp.Tfid, path.Path, *serr.Err) {
	path, err := serr.PathSplitErr(pn)
	if err != nil {
		return sp.NoFid, path, err
	}
	root, rp, rest := mfs.Root(path)
	db.DPrintf(db.MEMFSSRV, "rootFid: %q root %v rp %q rest %v\n", pn, root, rp, rest)
	fid, ok := mfs.roots.lookupAlloc(rp.String())
	if ok {
		db.DPrintf(db.MEMFSSRV, "rootFid: %q new fid %d\n", pn, fid)
		mfs.ps.NewRootFid(fid, mfs.ctx, root, rp)
	}
	return fid, rest, nil
}

// Returns FsObj for pn and the path from the root to FsObj
func (mfs *MemFs) lookupWalk(pn string) (fs.FsObj, path.Path, *serr.Err) {
	fid, path, err := mfs.rootFid(pn)
	if err != nil {
		return nil, nil, err
	}
	db.DPrintf(db.MEMFSSRV, "Create %q %v path %v\n", pn, fid, path)
	_, _, lo, err := mfs.ps.LookupWalk(fid, path, false, lockmap.RLOCK)
	if err != nil {
		db.DPrintf(db.MEMFSSRV, "LookupWalk %v err %v\n", path.Dir(), err)
		return nil, nil, err
	}
	return lo, path, err
}

func (mfs *MemFs) lookup(path path.Path, ltype lockmap.Tlock) (fs.FsObj, *lockmap.PathLock, *serr.Err) {
	d, _, path := mfs.Root(path)
	lk := mfs.plt.Acquire(mfs.ctx, rootP, ltype)
	if len(path) == 0 {
		return d, lk, nil
	}
	_, lo, lk, _, err := namei.Walk(mfs.plt, mfs.ctx, d, lk, rootP, path, nil, ltype)
	if err != nil {
		mfs.plt.Release(mfs.ctx, lk, ltype)
		return nil, nil, err
	}
	return lo, lk, nil
}

func (mfs *MemFs) lookupParent(path path.Path, ltype lockmap.Tlock) (fs.Dir, *lockmap.PathLock, *serr.Err) {
	lo, lk, err := mfs.lookup(path, ltype)
	if err != nil {
		return nil, nil, err
	}
	d := lo.(fs.Dir)
	return d, lk, nil
}

func (mfs *MemFs) NewDev(pn string, dev fs.FsObj) *serr.Err {
	db.DPrintf(db.MEMFSSRV, "NewDev %q\n", pn)
	path, err := serr.PathSplitErr(pn)
	if err != nil {
		return err
	}
	d, lk, err := mfs.lookupParent(path.Dir(), lockmap.WLOCK)
	db.DPrintf(db.MEMFSSRV, "lookupParent %q dir %v err %v\n", pn, d, err)
	if err != nil {
		return err
	}
	defer mfs.plt.Release(mfs.ctx, lk, lockmap.WLOCK)
	dev.SetParent(d)
	return dir.MkNod(mfs.ctx, d, path.Base(), dev)
}

func (mfs *MemFs) MkNod(pn string, i fs.FsObj) *serr.Err {
	db.DPrintf(db.MEMFSSRV, "MkNod %q\n", pn)
	path, err := serr.PathSplitErr(pn)
	if err != nil {
		return err
	}
	d, lk, err := mfs.lookupParent(path.Dir(), lockmap.WLOCK)
	db.DPrintf(db.MEMFSSRV, "MkNod: lookupParent %q dir %v err %v\n", pn, d, err)
	if err != nil {
		return err
	}
	defer mfs.plt.Release(mfs.ctx, lk, lockmap.WLOCK)
	return dir.MkNod(mfs.ctx, d, path.Base(), i)
}

func (mfs *MemFs) Create(pn string, p sp.Tperm, m sp.Tmode, lid sp.TleaseId) (fs.FsObj, *serr.Err) {
	lo, path, err := mfs.lookupWalk(filepath.Dir(pn))
	db.DPrintf(db.MEMFSSRV, "Create %q dir %v base %q\n", pn, lo, path.Base())
	_, nf, err := mfs.CreateObj(mfs.ctx, lo, path, filepath.Base(pn), p, m, lid, sp.NoFence())
	if err != nil {
		db.DPrintf(db.MEMFSSRV, "CreateObj %q %v err %v\n", pn, nf, err)
		return nil, err
	}
	return nf.Pobj().Obj(), nil
}

func (mfs *MemFs) Remove(pn string) *serr.Err {
	lo, path, err := mfs.lookupWalk(pn)
	if err != nil {
		return err
	}
	db.DPrintf(db.MEMFSSRV, "Remove %q %v %v\n", pn, lo, path.Base())
	return mfs.RemoveObj(mfs.ctx, lo, path, sp.NoFence())
}

func (mfs *MemFs) Open(pn string, m sp.Tmode) (fs.FsObj, *serr.Err) {
	lo, _, err := mfs.lookupWalk(pn)
	if err != nil {
		return nil, err
	}
	no, _, err := mfs.OpenObj(mfs.ctx, lo, m)
	return no, nil
}

func (mfs *MemFs) Dump() error {
	d, _, path := mfs.Root(rootP)
	s, err := d.(*dir.DirImpl).Dump()
	if err != nil {
		return err
	}
	db.DPrintf("MEMFSSRV", "Dump: %v %v", path, s)
	return nil
}
